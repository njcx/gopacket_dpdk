package dpdk

/*
#include "dpdk_wrapper.h"
#include "dpdk_bpf.h"
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

const (
	BURST_SIZE = 32
)

type DPDKHandle struct {
	portID      uint16
	bpfFilter   *C.dpdk_bpf_filter
	Initialized bool
	mbufs       []*C.struct_rte_mbuf
	currentIdx  int
	nbRx        int
	mu          sync.Mutex
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
