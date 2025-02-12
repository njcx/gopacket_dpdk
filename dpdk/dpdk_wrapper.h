#ifndef DPDK_WRAPPER_H
#define DPDK_WRAPPER_H

#include <rte_eal.h>
#include <rte_ethdev.h>
#include <rte_mbuf.h>


// 初始化DPDK环境
int init_dpdk(int argc, char **argv);

// 启动端口
int start_port(uint16_t port_id);


// 包装 rte_pktmbuf_mtod 宏的函数
void* get_mbuf_data(struct rte_mbuf* mbuf);

// 获取数据包长度
uint16_t get_mbuf_data_len(struct rte_mbuf* mbuf);

// 包装 rte_pktmbuf_free 的函数
void free_mbuf(struct rte_mbuf* mbuf);


#endif