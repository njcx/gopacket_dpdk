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
		fmt.Printf("源MAC: %s, 目标MAC: %s\n", eth.SrcMAC, eth.DstMAC)
	}

	// 解析IP层
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, ok := ipLayer.(*layers.IPv4)
		if !ok {
			fmt.Println("无法解析IPv4层")
			return
		}
		fmt.Printf("源IP: %s, 目标IP: %s\n", ip.SrcIP, ip.DstIP)
	}

	var resultDataList []map[string]string
	// 解析DNS层
	dnsLayer := packet.Layer(layers.LayerTypeDNS)
	if dnsLayer != nil {
		dns, ok := dnsLayer.(*layers.DNS)
		if !ok {
			fmt.Println("无法解析DNS层")
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
	// 初始化DPDK

	if os.Geteuid() != 0 {
		log.Fatal(" 需要root权限执行")
	}

	if err := dpdk.InitDPDK(); err != nil {
		log.Fatalf("初始化DPDK失败: %v", err)
	}

	handle, err := dpdk.NewDPDKHandle(0, "udp and port 53")
	if err != nil {
		log.Fatalf("创建DPDK处理器失败: %v", err)
	}
	defer handle.Close()

	// 开始接收和处理数据包
	handle.ReceivePackets(processPacket)
}
