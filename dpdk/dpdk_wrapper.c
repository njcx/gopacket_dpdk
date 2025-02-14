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
struct rte_mempool *mbuf_pool = NULL;

int init_dpdk(int argc, char **argv) {

    int ret;
    ret = rte_eal_init(argc, argv);
    printf("DPDK Version: %s\n", rte_version());
    if (ret < 0) {
        printf("Error: Cannot init EAL: %s\n", rte_strerror(rte_errno));
        return -1;
    }
     return ret;
}

int init_port(uint16_t port_id) {
    int ret;
    unsigned nb_ports;
    uint16_t i;
    struct rte_eth_conf port_conf = port_conf_default;

    nb_ports = rte_eth_dev_count_avail();
    printf("Number of available ports: %u\n", nb_ports);
    if (nb_ports < 1) {
        printf("Warning: No Ethernet ports available\n");
        return -1;
    }
    printf("Configuring port %u...\n", port_id);

    if (!rte_eth_dev_is_valid_port(port_id)) {
        printf("Invalid port ID %u\n", port_id);
        return -1;
    }

    mbuf_pool = rte_pktmbuf_pool_create("MBUF_POOL", NUM_MBUFS,
                                        MBUF_CACHE_SIZE, 0,
                                        RTE_MBUF_DEFAULT_BUF_SIZE,
                                        rte_socket_id());
    if (mbuf_pool == NULL) {
        printf("Error: Cannot create mbuf pool\n");
        return -1;
    }

    ret = rte_eth_dev_configure(port_id, 1, 1, &port_conf);
    if (ret < 0) {
        printf("Warning: Cannot configure port %u\n", port_id);
        return -1;
    }

    ret = rte_eth_rx_queue_setup(port_id, 0, RX_RING_SIZE,
                                 rte_eth_dev_socket_id(port_id),
                                 NULL, mbuf_pool);
    if (ret < 0) {
        printf("Warning: Cannot setup RX queue for port %u\n", port_id);
        return -1;
    }

    ret = rte_eth_tx_queue_setup(port_id, 0, TX_RING_SIZE,
                                 rte_eth_dev_socket_id(port_id),
                                 NULL);
    if (ret < 0) {
        printf("Warning: Cannot setup TX queue for port %u\n", port_id);
        return -1;
    }

    return 0;
}




int start_port(uint16_t port_id) {
    int ret = rte_eth_dev_start(port_id);
    if (ret < 0) {
        return ret;
    }
    rte_eth_promiscuous_enable(port_id);
    printf("Port %u started successfully\n", port_id);
    return 0;
}

void stop_port(uint16_t port_id) {
    if (rte_eth_dev_is_valid_port(port_id)) {
        printf("Stopping port %u...\n", port_id);
        rte_eth_dev_stop(port_id);
        rte_eth_dev_close(port_id);
    }
}

void cleanup_dpdk(void) {
    printf("Cleaning up DPDK resources...\n");
    if (mbuf_pool != NULL) {
        rte_mempool_free(mbuf_pool);
        mbuf_pool = NULL;
    }

    rte_eal_cleanup();
    printf("DPDK cleanup completed\n");
}

uint16_t receive_packets(uint16_t port_id, struct rte_mbuf **rx_pkts, uint16_t nb_pkts) {
    return rte_eth_rx_burst(port_id, 0, rx_pkts, nb_pkts);
}

uint16_t send_packets(uint16_t port_id, struct rte_mbuf **tx_pkts, uint16_t nb_pkts) {
    return rte_eth_tx_burst(port_id, 0, tx_pkts, nb_pkts);
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


uint16_t get_nb_ports(void) {
    return rte_eth_dev_count_avail();
}

int get_port_status(uint16_t port_id) {
    struct rte_eth_link link;
    int ret = rte_eth_link_get(port_id, &link);
    if (ret < 0) return ret;
    return link.link_status ? 1 : 0;
}

void print_port_info(uint16_t port_id) {
    struct rte_eth_dev_info dev_info;
    struct rte_eth_link link;
    struct rte_eth_stats stats;

    if (rte_eth_dev_info_get(port_id, &dev_info) != 0) {
        printf("Failed to get port info\n");
        return;
    }

    if (rte_eth_link_get(port_id, &link) != 0) {
        printf("Failed to get link info\n");
        return;
    }

    if (rte_eth_stats_get(port_id, &stats) != 0) {
        printf("Failed to get port statistics\n");
        return;
    }

    printf("\nPort %u information:\n", port_id);
    printf("Driver name: %s\n", dev_info.driver_name);
    printf("Link status: %s\n", link.link_status ? "up" : "down");
    printf("Link speed: %u Mbps\n", link.link_speed);
    printf("Link duplex: %s\n", link.link_duplex ? "full" : "half");
    printf("RX packets: %lu\n", stats.ipackets);
    printf("TX packets: %lu\n", stats.opackets);
    printf("RX bytes: %lu\n", stats.ibytes);
    printf("TX bytes: %lu\n", stats.obytes);
    printf("RX errors: %lu\n", stats.ierrors);
    printf("TX errors: %lu\n", stats.oerrors);
}