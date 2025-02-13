#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <signal.h>
#include <rte_eal.h>
#include <rte_ethdev.h>
#include <rte_mbuf.h>
#include <rte_version.h>

#define RX_RING_SIZE 512
#define TX_RING_SIZE 512
#define NUM_MBUFS 8191
#define MBUF_CACHE_SIZE 250
#define BURST_SIZE 32

static const struct rte_eth_conf port_conf_default = {
        .rxmode = {
                .max_lro_pkt_size = RTE_ETHER_MAX_LEN,
        },
};

static volatile int force_quit = 0;

void signal_handler(int signum);

int main(int argc, char *argv[])
{
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

        // Configure the Ethernet device
        ret = rte_eth_dev_configure(portid, 1, 1, &port_conf);
        if (ret < 0) {
            rte_exit(EXIT_FAILURE, "Cannot configure port %u\n", portid);
        }

        // RX queue setup
        ret = rte_eth_rx_queue_setup(portid, 0, RX_RING_SIZE,
                                     rte_eth_dev_socket_id(portid),
                                     NULL, mbuf_pool);
        if (ret < 0) {
            rte_exit(EXIT_FAILURE, "Cannot setup RX queue for port %u\n", portid);
        }

        // TX queue setup - using TX_RING_SIZE
        ret = rte_eth_tx_queue_setup(portid, 0, TX_RING_SIZE,
                                     rte_eth_dev_socket_id(portid),
                                     NULL);
        if (ret < 0) {
            rte_exit(EXIT_FAILURE, "Cannot setup TX queue for port %u\n", portid);
        }

        // Start the Ethernet port
        ret = rte_eth_dev_start(portid);
        if (ret < 0) {
            rte_exit(EXIT_FAILURE, "Cannot start port %u\n", portid);
        }

        printf("Port %u started successfully\n", portid);
    }

    printf("Running...\n");

    // 注册信号处理函数
    signal(SIGINT, signal_handler);
    signal(SIGTERM, signal_handler);

    // 接收数据包
    struct rte_mbuf *bufs[BURST_SIZE];
    while (!force_quit) {
        for (portid = 0; portid < nb_ports; portid++) {
            const uint16_t nb_rx = rte_eth_rx_burst(portid, 0, bufs, BURST_SIZE);
            if (nb_rx == 0)
                continue;

            for (i = 0; i < nb_rx && i < BURST_SIZE; i++) {
                if (bufs[i] == NULL) {
                    rte_exit(EXIT_FAILURE, "Received NULL packet buffer at index %u\n", i);
                }
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
        force_quit = 1;
    }
}