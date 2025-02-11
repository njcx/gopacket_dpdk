#include "dpdk_wrapper.h"

#if RTE_VERSION < RTE_VERSION_NUM(20, 0, 0, 0)
    #define RTE_ETH_MQ_RX_RSS ETH_MQ_RX_RSS
    #define RTE_ETH_RSS_IP ETH_RSS_IP
    #define RTE_ETH_RSS_UDP ETH_RSS_UDP
    #define RTE_ETH_RSS_TCP ETH_RSS_TCP
    #define RTE_ETH_MQ_TX_NONE ETH_MQ_TX_NONE
    #define RTE_ETHER_MAX_LEN ETHER_MAX_LEN
    #define RTE_ETH_DEV_NO_SCATTERED_RX DEV_RX_OFFLOAD_SCATTER
    #define RTE_ETH_RX_OFFLOAD_CHECKSUM DEV_RX_OFFLOAD_CHECKSUM
#endif

int init_dpdk(int argc, char **argv) {
    int ret = rte_eal_init(argc, argv);
    if (ret < 0) {
        return -1;
    }
    return ret;
}

int init_port(uint16_t port_id, uint16_t rx_rings, uint16_t tx_rings) {

   struct rte_eth_conf port_conf;
   memset(&port_conf, 0, sizeof(struct rte_eth_conf));

    port_conf.rxmode.max_rx_pkt_len = RTE_ETHER_MAX_LEN;
    port_conf.txmode.mq_mode = RTE_ETH_MQ_TX_NONE;
    #if RTE_VERSION >= RTE_VERSION_NUM(20, 0, 0, 0)
        port_conf.rxmode.offloads = RTE_ETH_RX_OFFLOAD_CHECKSUM;
    #else
        port_conf.rxmode.offloads = DEV_RX_OFFLOAD_CHECKSUM;
    #endif

    /* RSS 配置 */
    port_conf.rx_adv_conf.rss_conf.rss_key = NULL;

    port_conf.rx_adv_conf.rss_conf.rss_hf = RTE_ETH_RSS_IP |
                                           RTE_ETH_RSS_UDP |
                                           RTE_ETH_RSS_TCP;

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
                                 uint16_t data_room_size, unsigned socket_id) {
    return rte_pktmbuf_pool_create(name, n, cache_size, priv_size,
                                 data_room_size, socket_id);
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
