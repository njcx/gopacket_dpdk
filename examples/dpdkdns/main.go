package main

import (
	"fmt"
	"github.com/njcx/gopacket_dpdk"
	"github.com/njcx/gopacket_dpdk/dpdk"
	"github.com/njcx/gopacket_dpdk/layers"
	"log"
	"os"
)

func processPacket(data []byte) {
	packet := gopacket_dpdk.NewPacket(data, layers.LayerTypeEthernet, gopacket_dpdk.Default)
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		eth, _ := ethernetLayer.(*layers.Ethernet)
		fmt.Printf("Source MAC: %s, Destination MAC: %s\n", eth.SrcMAC, eth.DstMAC)
	}
	// Parse IP layer
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, ok := ipLayer.(*layers.IPv4)
		if !ok {
			fmt.Println("Failed to parse IPv4 layer")
			return
		}
		fmt.Printf("Source IP: %s, Destination IP: %s\n", ip.SrcIP, ip.DstIP)
	}

	var resultDataList []map[string]string
	// 解析DNS层
	dnsLayer := packet.Layer(layers.LayerTypeDNS)
	if dnsLayer != nil {
		dns, ok := dnsLayer.(*layers.DNS)
		if !ok {
			fmt.Println("Failed to parse DNS layer")
			return
		}
		if !dns.QR {
			for _, dnsQuestion := range dns.Questions {
				if len(dns.Questions) == 0 {
					continue
				}
				resultdata := make(map[string]string)
				resultdata["source"] = "dns"
				resultdata["domain"] = string(dnsQuestion.Name)
				resultdata["type"] = string(dnsQuestion.Type)
				resultdata["class"] = string(dnsQuestion.Class)
				resultDataList = append(resultDataList, resultdata)
			}
			for _, data := range resultDataList {
				fmt.Printf("%+v\n", data)
			}
		}
	}
}

func main() {
	// Initialize DPDK
	if os.Geteuid() != 0 {
		log.Fatal("Root permission is required to execute")
	}
	if err := dpdk.InitDPDK(); err != nil {
		log.Fatalf("Failed to initialize DPDK: %v", err)
	}
	handle, err := dpdk.NewDPDKHandle(0, "udp and port 53")
	if err != nil {
		log.Fatalf("Failed to create DPDK handler: %v", err)
	}
	defer handle.Close()
	// Start receiving and processing packets
	handle.ReceivePacketsCallBack(processPacket)
}
