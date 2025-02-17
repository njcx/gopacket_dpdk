package dpdk

/*
#include "dpdk_wrapper.h"
#include "dpdk_bpf.h"
*/
import "C"
import (
	"fmt"
	"sync"
	"time"
	"unsafe"
)

const (
	BURST_SIZE = 32
)

type BandwidthStats struct {
	RxBytesPerSecond   float64 // Received bandwidth (bytes/s)
	TxBytesPerSecond   float64 // Transmitted bandwidth (bytes/s)
	RxPacketsPerSecond float64 // Received packet rate (packets/s)
	TxPacketsPerSecond float64 // Transmitted packet rate (packets/s)
	Timestamp          time.Time
}

type PacketLossStats struct {
	RxDropped    uint64  // Number of received dropped packets
	TxDropped    uint64  // Number of transmitted dropped packets
	RxErrors     uint64  // Number of received errors
	TxErrors     uint64  // Number of transmitted errors
	RxLossRate   float64 // Received packet loss rate (percentage)
	TxLossRate   float64 // Transmitted packet loss rate (percentage)
	RxErrorRate  float64 // Received error rate (percentage)
	TxErrorRate  float64 // Transmitted error rate (percentage)
	TotalPackets uint64  // Total number of packets
	Timestamp    time.Time
}

type DPDKHandle struct {
	portID        uint16
	bpfFilter     *C.dpdk_bpf_filter
	Initialized   bool
	mbufs         []*C.struct_rte_mbuf
	currentIdx    int
	nbRx          int
	mu            sync.Mutex
	lastStats     *C.struct_rte_eth_stats
	lastStatsTime time.Time
}

func InitDPDK() error {
	args := []string{""}
	argc := C.int(len(args))
	cargs := make([]*C.char, len(args))
	for i, arg := range args {
		cargs[i] = C.CString(arg)
		defer C.free(unsafe.Pointer(cargs[i]))
	}
	ret := C.init_dpdk(argc, (**C.char)(&cargs[0]))
	if ret < 0 {
		return fmt.Errorf("DPDK initialization failed: %d", ret)
	}
	return nil

}

func NewDPDKHandle(portID uint16, bpfExpression string) (*DPDKHandle, error) {
	handle := &DPDKHandle{
		portID:    portID,
		mbufs:     make([]*C.struct_rte_mbuf, BURST_SIZE),
		bpfFilter: &C.dpdk_bpf_filter{},
	}
	if bpfExpression != "" {
		cExpr := C.CString(bpfExpression)
		defer C.free(unsafe.Pointer(cExpr))
		if ret := C.init_bpf_filter(handle.bpfFilter, cExpr, 0xffffff00); ret != 0 {
			return nil, fmt.Errorf("BPF filter initialization failed: %d", ret)
		}
	}
	if ret := C.init_port(C.uint16_t(portID)); ret != 0 {
		return nil, fmt.Errorf("port initialization failed: %d", ret)
	}
	if ret := C.start_port(C.uint16_t(portID)); ret != 0 {
		return nil, fmt.Errorf("port start failed: %d", ret)
	}
	handle.Initialized = true
	return handle, nil
}

func (h *DPDKHandle) Close() {
	if h.bpfFilter != nil {
		C.free_bpf_filter(h.bpfFilter)
	}
	if h.Initialized {
		C.stop_port(C.uint16_t(h.portID))
		C.cleanup_dpdk()
		h.Initialized = false
	}

}

func (h *DPDKHandle) ReceivePacketsCallBack(callback func([]byte)) {
	burstSize := C.uint16_t(BURST_SIZE)
	mbufs := make([]*C.struct_rte_mbuf, BURST_SIZE)

	for {
		nb_rx := C.rte_eth_rx_burst(C.uint16_t(h.portID), 0,
			(**C.struct_rte_mbuf)(unsafe.Pointer(&mbufs[0])),
			burstSize)
		for i := 0; i < int(nb_rx); i++ {
			mbuf := mbufs[i]
			data := C.get_mbuf_data(mbuf)
			length := C.get_mbuf_data_len(mbuf)

			if h.bpfFilter != nil {
				if C.apply_bpf_filter(h.bpfFilter,
					(*C.uchar)(unsafe.Pointer(data)),
					C.uint32_t(length)) == 0 {
					C.free_mbuf(mbuf)
					continue
				}
			}
			goData := C.GoBytes(unsafe.Pointer(data), C.int(length))
			callback(goData)
			C.free_mbuf(mbuf)
		}
	}
}

func (h *DPDKHandle) ReadPacket() ([]byte, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.currentIdx >= h.nbRx {
		h.nbRx = int(C.receive_packets(C.uint16_t(h.portID),
			(**C.struct_rte_mbuf)(unsafe.Pointer(&h.mbufs[0])),
			C.uint16_t(BURST_SIZE)))

		h.currentIdx = 0

		if h.nbRx == 0 {
			return nil, nil
		}
	}

	if h.currentIdx < 0 || h.currentIdx >= len(h.mbufs) {
		return nil, fmt.Errorf("currentIdx out of bounds: %d", h.currentIdx)
	}

	mbuf := h.mbufs[h.currentIdx]
	data := C.get_mbuf_data(mbuf)
	length := C.get_mbuf_data_len(mbuf)

	packet := C.GoBytes(unsafe.Pointer(data), C.int(length))
	C.free_mbuf(mbuf)
	h.currentIdx++

	return packet, nil
}

func (h *DPDKHandle) SendPackets(packets [][]byte) (uint16, error) {
	if len(packets) == 0 {
		return 0, nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	mbufs := make([]*C.struct_rte_mbuf, len(packets))

	sent := C.send_packets(C.uint16_t(h.portID),
		(**C.struct_rte_mbuf)(unsafe.Pointer(&mbufs[0])),
		C.uint16_t(len(packets)))

	return uint16(sent), nil
}

func (h *DPDKHandle) IsPortUp() bool {
	status := C.get_port_status(C.uint16_t(h.portID))
	return status > 0
}

func (h *DPDKHandle) PrintInfo() {
	C.print_port_info(C.uint16_t(h.portID))
}

func (h *DPDKHandle) GetBandwidth() (*BandwidthStats, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	var currentStats C.struct_rte_eth_stats
	if ret := C.rte_eth_stats_get(C.uint16_t(h.portID), &currentStats); ret != 0 {
		return nil, fmt.Errorf("failed to get port statistics: %d", ret)
	}

	currentTime := time.Now()

	if h.lastStats == nil {
		h.lastStats = &C.struct_rte_eth_stats{}
		*h.lastStats = currentStats
		h.lastStatsTime = currentTime
		return &BandwidthStats{
			Timestamp: currentTime,
		}, nil
	}

	duration := currentTime.Sub(h.lastStatsTime).Seconds()
	if duration == 0 {
		return nil, fmt.Errorf("too frequent calls to GetBandwidth")
	}

	stats := &BandwidthStats{
		RxBytesPerSecond:   float64(currentStats.ibytes-h.lastStats.ibytes) / duration,
		TxBytesPerSecond:   float64(currentStats.obytes-h.lastStats.obytes) / duration,
		RxPacketsPerSecond: float64(currentStats.ipackets-h.lastStats.ipackets) / duration,
		TxPacketsPerSecond: float64(currentStats.opackets-h.lastStats.opackets) / duration,
		Timestamp:          currentTime,
	}

	*h.lastStats = currentStats
	h.lastStatsTime = currentTime

	return stats, nil
}

func (h *DPDKHandle) PrintBandwidth() error {
	stats, err := h.GetBandwidth()
	if err != nil {
		return err
	}

	fmt.Printf("\nBandwidth Statistics (at %s):\n", stats.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("RX: %.2f Mbps (%.2f packets/s)\n", stats.RxBytesPerSecond*8/1000000, stats.RxPacketsPerSecond)
	fmt.Printf("TX: %.2f Mbps (%.2f packets/s)\n", stats.TxBytesPerSecond*8/1000000, stats.TxPacketsPerSecond)

	return nil
}

func (h *DPDKHandle) GetPacketLoss() (*PacketLossStats, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	var currentStats C.struct_rte_eth_stats
	if ret := C.rte_eth_stats_get(C.uint16_t(h.portID), &currentStats); ret != 0 {
		return nil, fmt.Errorf("failed to get port statistics: %d", ret)
	}

	currentTime := time.Now()

	if h.lastStats == nil {
		h.lastStats = &C.struct_rte_eth_stats{}
		*h.lastStats = currentStats
		h.lastStatsTime = currentTime
		return &PacketLossStats{
			Timestamp: currentTime,
		}, nil
	}

	totalRxPackets := uint64(currentStats.ipackets)
	totalTxPackets := uint64(currentStats.opackets)
	rxDropped := uint64(currentStats.imissed + currentStats.rx_nombuf)
	txDropped := uint64(currentStats.oerrors)
	rxErrors := uint64(currentStats.ierrors)
	txErrors := uint64(currentStats.oerrors)

	totalRxAttempted := totalRxPackets + rxDropped + rxErrors
	totalTxAttempted := totalTxPackets + txDropped + txErrors

	rxLossRate := float64(0)
	txLossRate := float64(0)
	rxErrorRate := float64(0)
	txErrorRate := float64(0)

	if totalRxAttempted > 0 {
		rxLossRate = float64(rxDropped) / float64(totalRxAttempted) * 100
		rxErrorRate = float64(rxErrors) / float64(totalRxAttempted) * 100
	}
	if totalTxAttempted > 0 {
		txLossRate = float64(txDropped) / float64(totalTxAttempted) * 100
		txErrorRate = float64(txErrors) / float64(totalTxAttempted) * 100
	}

	stats := &PacketLossStats{
		RxDropped:    rxDropped,
		TxDropped:    txDropped,
		RxErrors:     rxErrors,
		TxErrors:     txErrors,
		RxLossRate:   rxLossRate,
		TxLossRate:   txLossRate,
		RxErrorRate:  rxErrorRate,
		TxErrorRate:  txErrorRate,
		TotalPackets: totalRxPackets + totalTxPackets,
		Timestamp:    currentTime,
	}

	*h.lastStats = currentStats
	h.lastStatsTime = currentTime

	return stats, nil
}

func (h *DPDKHandle) PrintPacketLoss() error {
	stats, err := h.GetPacketLoss()
	if err != nil {
		return err
	}

	fmt.Printf("\nPacket Loss Statistics (at %s):\n", stats.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("Total Packets Processed: %d\n", stats.TotalPackets)
	fmt.Printf("RX Dropped: %d (%.2f%% of total received)\n", stats.RxDropped, stats.RxLossRate)
	fmt.Printf("TX Dropped: %d (%.2f%% of total transmitted)\n", stats.TxDropped, stats.TxLossRate)
	fmt.Printf("RX Errors: %d (%.2f%% of total received)\n", stats.RxErrors, stats.RxErrorRate)
	fmt.Printf("TX Errors: %d (%.2f%% of total transmitted)\n", stats.TxErrors, stats.TxErrorRate)

	return nil
}