// dpdk_bpf.h
#ifndef DPDK_BPF_H
#define DPDK_BPF_H

#include <pcap/pcap.h>
#include <pcap/bpf.h>

// BPF过滤器结构
typedef struct {
    struct bpf_program prog;    // 编译后的BPF程序
    int optimized;             // 是否已优化
} dpdk_bpf_filter;

// 初始化BPF过滤器
int init_bpf_filter(dpdk_bpf_filter *filter, const char *expression,
                   uint32_t netmask);

// 应用BPF过滤器到数据包
int apply_bpf_filter(dpdk_bpf_filter *filter, const unsigned char *packet,
                    uint32_t len);

// 释放BPF过滤器
void free_bpf_filter(dpdk_bpf_filter *filter);

#endif