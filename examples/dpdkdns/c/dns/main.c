#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <signal.h>
#include <rte_eal.h>
#include <rte_ethdev.h>
#include <rte_mbuf.h>
#include <rte_ether.h>
#include <rte_version.h>
#include <rte_ip.h>
#include <rte_udp.h>

#define RX_RING_SIZE 512
#define TX_RING_SIZE 512
#define NUM_MBUFS 8191
#define MBUF_CACHE_SIZE 250
#define BURST_SIZE 32

// DNS header structure
struct dns_hdr {
    uint16_t id;       // identification number
    uint16_t flags;    // DNS flags
    uint16_t qdcount;  // number of question entries
    uint16_t ancount;  // number of answer entries
    uint16_t nscount;  // number of authority entries
    uint16_t arcount;  // number of resource entries
} __attribute__((__packed__));

static const struct rte_eth_conf port_conf_default = {
        .rxmode = {
                .max_lro_pkt_size = RTE_ETHER_MAX_LEN,
        },
};

static volatile int force_quit = 0;

void signal_handler(int signum);

// 函数用于解析DNS查询中的域名
static void print_dns_name(const unsigned char *reader, const unsigned char *buffer, int *count) {
    unsigned char name[256];
    unsigned int p = 0;
    unsigned int jumped = 0;
    unsigned int offset;
    int i, j;

    *count = 1;
    name[0] = '\0';

    while(*reader != 0) {
        if(*reader >= 192) {  // Compression used
            offset = (*reader)*256 + *(reader+1) - 49152;
            reader = buffer + offset - 1;
            if(jumped == 0)
                *count += 1;
            jumped = 1;
        }
        else {
            name[p++] = *reader;
        }
        reader = reader + 1;
        if(jumped == 0)
            *count += 1;
    }

    name[p] = '\0';

    // 现在转换成可读格式
    for(i = 0; i < (int)strlen((const char*)name); i++) {
        p = name[i];
        for(j = 0; j < (int)p; j++) {
            printf("%c", name[i+1+j]);
        }
        i += p;
        if(i < (int)strlen((const char*)name) - 1)
            printf(".");
    }
    printf("\n");
}

// 处理DNS数据包的函数
static void process_dns_packet(struct rte_mbuf *mbuf) {
    struct rte_ether_hdr *eth_hdr;
    struct rte_ipv4_hdr *ip_hdr;
    struct rte_udp_hdr *udp_hdr;
    struct dns_hdr *dns_hdr;
    const unsigned char *reader;

    // 获取以太网头部
    eth_hdr = rte_pktmbuf_mtod(mbuf, struct rte_ether_hdr *);

    // 检查是否是IPv4数据包
    if (rte_be_to_cpu_16(eth_hdr->ether_type) != RTE_ETHER_TYPE_IPV4) {
        return;
    }

    // 获取IP头部
    ip_hdr = (struct rte_ipv4_hdr *)(eth_hdr + 1);

    // 检查是否是UDP数据包
    if (ip_hdr->next_proto_id != IPPROTO_UDP) {
        return;
    }

    // 获取UDP头部
    udp_hdr = (struct rte_udp_hdr *)((unsigned char *)ip_hdr + (ip_hdr->version_ihl & 0x0F) * 4);

    // 检查是否是DNS数据包（DNS默认端口为53）
    if (!(rte_be_to_cpu_16(udp_hdr->dst_port) == 53 || rte_be_to_cpu_16(udp_hdr->src_port) == 53)) {
        return;
    }

    // 获取DNS头部
    dns_hdr = (struct dns_hdr *)((unsigned char *)udp_hdr + sizeof(struct rte_udp_hdr));

    // 将读指针移到查询部分
    reader = (unsigned char *)(dns_hdr + 1);

    printf("DNS Query: ");
    int count = 0;
    print_dns_name(reader, (unsigned char *)dns_hdr, &count);
}

int main(int argc, char *argv[]) {
    int ret;
    unsigned nb_ports;
    uint16_t portid;
    uint16_t i;
    struct rte_mempool *mbuf_pool = NULL;
    struct rte_eth_conf port_conf = port_conf_default;

    // 初始化DPDK环境
    ret = rte_eal_init(argc, argv);
    if (ret < 0) {
        rte_exit(EXIT_FAILURE, "Cannot init EAL\n");
    }
    printf("DPDK Version: %s\n", rte_version());
    nb_ports = rte_eth_dev_count_avail();
    printf("Number of available ports: %u\n", nb_ports);
    if (nb_ports < 1) {
        rte_exit(EXIT_FAILURE, "No Ethernet ports\n");
    }

    // 分配内存池
    mbuf_pool = rte_pktmbuf_pool_create("MBUF_POOL", NUM_MBUFS,
                                        MBUF_CACHE_SIZE, 0,
                                        RTE_MBUF_DEFAULT_BUF_SIZE,
                                        rte_socket_id());
    if (mbuf_pool == NULL) {
        rte_exit(EXIT_FAILURE, "Cannot create mbuf pool\n");
    }

    // 配置并启动网卡
    for (portid = 0; portid < nb_ports; portid++) {
        printf("Configuring port %u...\n", portid);

        ret = rte_eth_dev_configure(portid, 1, 1, &port_conf);
        if (ret < 0) {
            rte_exit(EXIT_FAILURE, "Cannot configure port %u\n", portid);
        }

        ret = rte_eth_rx_queue_setup(portid, 0, RX_RING_SIZE,
                                     rte_eth_dev_socket_id(portid),
                                     NULL, mbuf_pool);
        if (ret < 0) {
            rte_exit(EXIT_FAILURE, "Cannot setup RX queue for port %u\n", portid);
        }

        ret = rte_eth_tx_queue_setup(portid, 0, TX_RING_SIZE,
                                     rte_eth_dev_socket_id(portid),
                                     NULL);
        if (ret < 0) {
            rte_exit(EXIT_FAILURE, "Cannot setup TX queue for port %u\n", portid);
        }

        ret = rte_eth_dev_start(portid);
        if (ret < 0) {
            rte_exit(EXIT_FAILURE, "Cannot start port %u\n", portid);
        }

        // 启用混杂模式以捕获所有DNS数据包
        rte_eth_promiscuous_enable(portid);

        printf("Port %u started successfully\n", portid);
    }

    printf("DNS monitoring started...\n");

    // 注册信号处理函数
    signal(SIGINT, signal_handler);
    signal(SIGTERM, signal_handler);

    // 接收和处理数据包
    struct rte_mbuf *bufs[BURST_SIZE];
    while (!force_quit) {
        for (portid = 0; portid < nb_ports; portid++) {
            const uint16_t nb_rx = rte_eth_rx_burst(portid, 0, bufs, BURST_SIZE);
            if (nb_rx == 0)
                continue;

            for (i = 0; i < nb_rx; i++) {
                process_dns_packet(bufs[i]);  // 处理DNS数据包
                rte_pktmbuf_free(bufs[i]);
            }
        }
    }

    // 停止网卡设备
    for (portid = 0; portid < nb_ports; portid++) {
        rte_eth_dev_stop(portid);
        rte_eth_dev_close(portid);
    }

    // 释放内存池
    if (mbuf_pool != NULL) {
        rte_mempool_free(mbuf_pool);
    }

    return 0;
}

void signal_handler(int signum) {
    if (signum == SIGINT || signum == SIGTERM) {
        printf("\nSignal %d received, preparing to exit...\n", signum);
        force_quit = 1;
    }
}