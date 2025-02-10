package main

import (
	"fmt"
	"github.com/njcx/gopacket"
	"github.com/njcx/gopacket/dpdk"
	"github.com/njcx/gopacket/layers"
	"log"
)

func processPacket(data []byte) {
	packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)

	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		eth, _ := ethernetLayer.(*layers.Ethernet)
		fmt.Printf("源MAC: %s, 目标MAC: %s\n", eth.SrcMAC, eth.DstMAC)
	}

	// 解析IP层
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		fmt.Printf("源IP: %s, 目标IP: %s\n", ip.SrcIP, ip.DstIP)
	}
	var resultdata = make(map[string]string)
	// 解析DNS层
	dnsLayer := packet.Layer(layers.LayerTypeDNS)
	if dnsLayer != nil {
		dns, _ := dnsLayer.(*layers.LayerTypeDNS)
		if !dns.QR {
			for _, dnsQuestion := range dns.Questions {
				resultdata["source"] = "dns"
				resultdata["src"] = dnsn.SrcIP
				resultdata["dst"] = dns.DstIP
				resultdata["domain"] = string(dnsQuestion.Name)
				resultdata["type"] = dnsQuestion.Type.String()
				resultdata["class"] = dnsQuestion.Class.String()
				fmt.Println(resultdata)

			}
		}

	}
}

func main() {
	// 初始化DPDK
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
