#ifndef DPDK_WRAPPER_H
#define DPDK_WRAPPER_H

#include <dpdk/rte_eal.h>
#include <dpdk/rte_ethdev.h>
#include <dpdk/rte_mbuf.h>

// 初始化DPDK环境
int init_dpdk(int argc, char **argv);

// 初始化端口
int init_port(uint16_t port_id, uint16_t rx_rings, uint16_t tx_rings);

// 启动端口
int start_port(uint16_t port_id);

// 创建内存池
struct rte_mempool* create_mempool(const char* name, unsigned n,
                                 unsigned cache_size, uint16_t priv_size,
                                 uint16_t data_room_size, int socket_id);

#endif