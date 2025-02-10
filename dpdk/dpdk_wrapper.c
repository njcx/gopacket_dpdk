#include "dpdk_wrapper.h"

int init_dpdk(int argc, char **argv) {
    int ret = rte_eal_init(argc, argv);
    if (ret < 0) {
        return -1;
    }
    return ret;
}

int init_port(uint16_t port_id, uint16_t rx_rings, uint16_t tx_rings) {
    struct rte_eth_conf port_conf = {
        .rxmode = {
            .max_rx_pkt_len = RTE_ETHER_MAX_LEN,
            .split_hdr_size = 0,
        },
        .txmode = {
            .mq_mode = ETH_MQ_TX_NONE,
        },
    };

    // 配置设备
    int ret = rte_eth_dev_configure(port_id, rx_rings, tx_rings, &port_conf);
    if (ret != 0) {
        return ret;
    }

    // 配置每个接收队列
    for (int i = 0; i < rx_rings; i++) {
        ret = rte_eth_rx_queue_setup(port_id, i, 1024,
                                   rte_eth_dev_socket_id(port_id),
                                   NULL, rte_pktmbuf_pool_create("mbuf_pool", 8192,
                                   256, 0, RTE_MBUF_DEFAULT_BUF_SIZE,
                                   rte_socket_id()));
        if (ret < 0) {
            return ret;
        }
    }

    // 配置每个发送队列
    for (int i = 0; i < tx_rings; i++) {
        ret = rte_eth_tx_queue_setup(port_id, i, 1024,
                                   rte_eth_dev_socket_id(port_id),
                                   NULL);
        if (ret < 0) {
            return ret;
        }
    }

    return 0;
}

int start_port(uint16_t port_id) {
    int ret = rte_eth_dev_start(port_id);
    if (ret < 0) {
        return ret;
    }

    // 设置混杂模式
    rte_eth_promiscuous_enable(port_id);
    return 0;
}

struct rte_mempool* create_mempool(const char* name, unsigned n,
                                 unsigned cache_size, uint16_t priv_size,
                                 uint16_t data_room_size, int socket_id) {
    return rte_pktmbuf_pool_create(name, n, cache_size, priv_size,
                                 data_room_size, socket_id);
}