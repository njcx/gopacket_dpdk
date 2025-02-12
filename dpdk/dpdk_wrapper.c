#include "dpdk_wrapper.h"

#define RX_RING_SIZE 1024
#define TX_RING_SIZE 1024
#define NUM_MBUFS 8191
#define MBUF_CACHE_SIZE 250
#define BURST_SIZE 32

static const struct rte_eth_conf port_conf_default = {
        .rxmode = {
                .max_lro_pkt_size = RTE_ETHER_MAX_LEN,
        },
};

int init_dpdk(int argc, char **argv) {

    int ret;
    unsigned nb_ports;
    uint16_t portid;
    uint16_t i;
    struct rte_mempool *mbuf_pool = NULL;
    struct rte_eth_conf port_conf = port_conf_default;

    // 初始化DPDK环境
    ret = rte_eal_init(argc, argv);
    if (ret < 0) {
        printf("Error: Cannot init EAL\n");
        return -1;
    }

    nb_ports = rte_eth_dev_count_avail();
    if (nb_ports < 1) {
        printf("Warning: No Ethernet ports available\n");
        return -1;
    }

    // 分配内存池
    mbuf_pool = rte_pktmbuf_pool_create("MBUF_POOL", NUM_MBUFS,
                                        MBUF_CACHE_SIZE, 0,
                                        RTE_MBUF_DEFAULT_BUF_SIZE,
                                        rte_socket_id());
    if (mbuf_pool == NULL) {
        printf("Error: Cannot create mbuf pool\n");
        return -1;
    }

    portid = 0;
    ret = rte_eth_dev_configure(portid, 1, 1, &port_conf);
    if (ret < 0) {
        printf("Warning: Cannot configure port %u\n", portid);
        return -1;
    }

    ret = rte_eth_rx_queue_setup(portid, 0, RX_RING_SIZE,
                                 rte_eth_dev_socket_id(portid),
                                 NULL, mbuf_pool);
    if (ret < 0) {
        printf("Warning: Cannot setup RX queue for port %u\n", portid);
        return -1;
    }

    ret = rte_eth_tx_queue_setup(portid, 0, TX_RING_SIZE,
                                 rte_eth_dev_socket_id(portid),
                                 NULL);
    if (ret < 0) {
        printf("Warning: Cannot setup TX queue for port %u\n", portid);
        return -1;
    }
}


int start_port(uint16_t port_id) {
    int ret = rte_eth_dev_start(port_id);
    if (ret < 0) {
        return ret;
    }

    rte_eth_promiscuous_enable(port_id);
    return 0;
}



void* get_mbuf_data(struct rte_mbuf* mbuf) {
    return rte_pktmbuf_mtod(mbuf, void*);
}

uint16_t get_mbuf_data_len(struct rte_mbuf* mbuf) {
    return rte_pktmbuf_data_len(mbuf);
}

void free_mbuf(struct rte_mbuf* mbuf) {
    rte_pktmbuf_free(mbuf);
}
