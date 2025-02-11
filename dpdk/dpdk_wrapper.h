#ifndef DPDK_WRAPPER_H
#define DPDK_WRAPPER_H

#include <rte_eal.h>
#include <rte_ethdev.h>
#include <rte_mbuf.h>
#include <rte_common.h>
#include <rte_version.h>
#include <rte_ether.h>


// 初始化DPDK环境
int init_dpdk(int argc, char **argv);

// 初始化端口
int init_port(uint16_t port_id, uint16_t rx_rings, uint16_t tx_rings);

// 启动端口
int start_port(uint16_t port_id);

// 创建内存池
struct rte_mempool* create_mempool(const char* name, unsigned n,
                                 unsigned cache_size, uint16_t priv_size,
                                 uint16_t data_room_size, unsigned socket_id);

// 包装 rte_pktmbuf_mtod 宏的函数
void* get_mbuf_data(struct rte_mbuf* mbuf);

// 获取数据包长度
uint16_t get_mbuf_data_len(struct rte_mbuf* mbuf);

// 包装 rte_pktmbuf_free 的函数
void free_mbuf(struct rte_mbuf* mbuf);


#endif