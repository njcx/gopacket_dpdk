// dpdk_bpf.h
#ifndef DPDK_BPF_H
#define DPDK_BPF_H

#include <stdio.h>
#include <string.h>
#include <stdint.h>
#include <pcap/pcap.h>
#include <pcap/bpf.h>


typedef struct {
    struct bpf_program prog;
    int optimized;
} dpdk_bpf_filter;


int init_bpf_filter(dpdk_bpf_filter *filter, const char *expression,
                   uint32_t netmask);


int apply_bpf_filter(dpdk_bpf_filter *filter, const unsigned char *packet,
                    uint32_t len);

void free_bpf_filter(dpdk_bpf_filter *filter);

#endif