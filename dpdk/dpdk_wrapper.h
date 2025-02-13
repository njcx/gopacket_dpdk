#ifndef DPDK_WRAPPER_H
#define DPDK_WRAPPER_H

#include <rte_eal.h>
#include <rte_ethdev.h>
#include <rte_mbuf.h>
#include <rte_version.h>


extern struct rte_mempool *mbuf_pool;

int init_dpdk(int argc, char **argv);

int init_port(uint16_t port_id);

int start_port(uint16_t port_id);

uint16_t get_nb_ports(void);

void stop_port(uint16_t port_id);

// 清理DPDK资源
void cleanup_dpdk(void);

// 接收数据包
uint16_t receive_packets(uint16_t port_id, struct rte_mbuf **rx_pkts, uint16_t nb_pkts);

// 发送数据包
uint16_t send_packets(uint16_t port_id, struct rte_mbuf **tx_pkts, uint16_t nb_pkts);

// 获取端口状态
int get_port_status(uint16_t port_id);

// 打印端口信息
void print_port_info(uint16_t port_id);


void* get_mbuf_data(struct rte_mbuf* mbuf);

uint16_t get_mbuf_data_len(struct rte_mbuf* mbuf);

void free_mbuf(struct rte_mbuf* mbuf);


#endif