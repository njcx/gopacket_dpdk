package dpdk

/*
#cgo CFLAGS: -m64 -pthread -O3 -march=native
#cgo LDFLAGS:  -lrte_eal -lrte_ethdev -lrte_mbuf -lrte_mempool  -lpcap
#include "dpdk_wrapper.h"
#include "dpdk_bpf.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

const (
	RX_RING_SIZE    = 1024
	TX_RING_SIZE    = 1024
	NUM_MBUFS       = 8191
	MBUF_CACHE_SIZE = 250
	BURST_SIZE      = 32
)

func InitDPDK() error {
	args := []string{"", "-l", "0-3", "-n", "4", "--proc-type=auto"}
	argc := C.int(len(args))

	// 将Go字符串转换为C字符串数组
	cargs := make([]*C.char, len(args))
	for i, arg := range args {
		cargs[i] = C.CString(arg)
		defer C.free(unsafe.Pointer(cargs[i]))
	}

	ret := C.init_dpdk(argc, (**C.char)(&cargs[0]))
	if ret < 0 {
		return fmt.Errorf("DPDK初始化失败: %d", ret)
	}
	return nil
}

type DPDKHandle struct {
	portID    uint16
	mempool   *C.struct_rte_mempool
	bpfFilter *C.dpdk_bpf_filter
}

// 创建新的DPDK处理器，包含BPF过滤器
func NewDPDKHandle(portID uint16, bpfExpression string) (*DPDKHandle, error) {
	handle := &DPDKHandle{
		portID:    portID,
		bpfFilter: &C.dpdk_bpf_filter{},
	}

	// 创建内存池
	mempoolName := C.CString(fmt.Sprintf("mbuf_pool_%d", portID))
	defer C.free(unsafe.Pointer(mempoolName))

	handle.mempool = C.create_mempool(mempoolName, NUM_MBUFS, MBUF_CACHE_SIZE,
		0, C.RTE_MBUF_DEFAULT_BUF_SIZE, C.rte_socket_id())
	if handle.mempool == nil {
		return nil, fmt.Errorf("创建内存池失败")
	}

	// 初始化BPF过滤器
	if bpfExpression != "" {
		cExpr := C.CString(bpfExpression)
		defer C.free(unsafe.Pointer(cExpr))

		if ret := C.init_bpf_filter(handle.bpfFilter, cExpr, 0xffffff00); ret != 0 {
			return nil, fmt.Errorf("BPF过滤器初始化失败: %d", ret)
		}
	}

	// 初始化端口
	if ret := C.init_port(C.uint16_t(portID), RX_RING_SIZE, TX_RING_SIZE); ret != 0 {
		return nil, fmt.Errorf("端口初始化失败: %d", ret)
	}

	// 启动端口
	if ret := C.start_port(C.uint16_t(portID)); ret != 0 {
		return nil, fmt.Errorf("端口启动失败: %d", ret)
	}

	return handle, nil
}

// 清理资源
func (h *DPDKHandle) Close() {
	if h.bpfFilter != nil {
		C.free_bpf_filter(h.bpfFilter)
	}
}

// 接收和过滤数据包
func (h *DPDKHandle) ReceivePackets(callback func([]byte)) {
	burstSize := C.uint16_t(BURST_SIZE)
	mbufs := make([]*C.struct_rte_mbuf, BURST_SIZE)

	for {
		// 接收一批数据包
		nb_rx := C.rte_eth_rx_burst(C.uint16_t(h.portID), 0,
			(**C.struct_rte_mbuf)(unsafe.Pointer(&mbufs[0])),
			burstSize)

		// 处理接收到的数据包
		for i := 0; i < int(nb_rx); i++ {
			mbuf := mbufs[i]
			data := C.rte_pktmbuf_mtod(mbuf, *C.char)
			length := C.rte_pktmbuf_data_len(mbuf)

			// 应用BPF过滤器
			if h.bpfFilter != nil {
				if C.apply_bpf_filter(h.bpfFilter,
					(*C.uchar)(unsafe.Pointer(data)),
					C.uint32_t(length)) == 0 {
					// 包不匹配过滤器，跳过处理
					C.rte_pktmbuf_free(mbuf)
					continue
				}
			}

			// 转换为Go切片
			goData := C.GoBytes(unsafe.Pointer(data), C.int(length))

			// 使用回调处理数据包
			callback(goData)

			// 释放mbuf
			C.rte_pktmbuf_free(mbuf)
		}
	}
}
