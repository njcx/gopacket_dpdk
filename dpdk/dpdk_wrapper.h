#ifndef DPDK_WRAPPER_H
#define DPDK_WRAPPER_H

#include <rte_eal.h>
#include <rte_ethdev.h>
#include <rte_mbuf.h>


int init_dpdk(int argc, char **argv);

int init_port(uint16_t port_id);

int start_port(uint16_t port_id);

void* get_mbuf_data(struct rte_mbuf* mbuf);

uint16_t get_mbuf_data_len(struct rte_mbuf* mbuf);

void free_mbuf(struct rte_mbuf* mbuf);


#endif