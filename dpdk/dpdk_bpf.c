// dpdk_bpf.c
#include "dpdk_bpf.h"
#include <stdio.h>
#include <string.h>

int init_bpf_filter(dpdk_bpf_filter *filter, const char *expression,
                   uint32_t netmask) {
    pcap_t *p;
    char errbuf[PCAP_ERRBUF_SIZE];

    // 创建临时pcap句柄用于编译BPF表达式
    p = pcap_open_dead(DLT_EN10MB, 65535);
    if (p == NULL) {
        return -1;
    }

    // 编译过滤器表达式
    if (pcap_compile(p, &filter->prog, expression, 1, netmask) < 0) {
        pcap_close(p);
        return -1;
    }

    // 优化BPF程序
    if (pcap_optimize(&filter->prog) < 0) {
        filter->optimized = 0;
    } else {
        filter->optimized = 1;
    }

    pcap_close(p);
    return 0;
}

int apply_bpf_filter(dpdk_bpf_filter *filter, const unsigned char *packet,
                    uint32_t len) {
    struct bpf_insn *insns = filter->prog.bf_insns;
    unsigned int res;

    if (!filter->optimized) {
        // 如果未优化，返回匹配以避免丢包
        return 1;
    }

    res = bpf_filter(insns, (uint8_t *)packet, len, len);
    return res != 0;
}

void free_bpf_filter(dpdk_bpf_filter *filter) {
    pcap_freecode(&filter->prog);
}