// dpdk_bpf.c
#include "dpdk_bpf.h"

int init_bpf_filter(dpdk_bpf_filter *filter, const char *expression,
                   uint32_t netmask) {
    pcap_t *p;
    char errbuf[PCAP_ERRBUF_SIZE];

    p = pcap_open_dead(DLT_EN10MB, 65535);
    if (p == NULL) {
        return -1;
    }

    if (pcap_compile(p, &filter->prog, expression, 1, netmask) < 0) {
        pcap_close(p);
        return -1;
    }

    filter->optimized = 1;
    pcap_close(p);
    return 0;
}

int apply_bpf_filter(dpdk_bpf_filter *filter, const unsigned char *packet,
                    uint32_t len) {
    struct bpf_insn *insns = filter->prog.bf_insns;
    unsigned int res;

    if (!filter->optimized) {
        return 1;
    }

    res = bpf_filter(insns, (uint8_t *)packet, len, len);
    return res != 0;
}

void free_bpf_filter(dpdk_bpf_filter *filter) {
    pcap_freecode(&filter->prog);
}