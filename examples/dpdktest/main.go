package main

import (
	"fmt"
	"github.com/njcx/gopacket_dpdk"
	"github.com/njcx/gopacket_dpdk/dpdk"
	"github.com/njcx/gopacket_dpdk/layers"
	"time"
)

func main() {

	if err := dpdk.InitDPDK(); err != nil {
		panic(fmt.Sprintf("Failed to initialize DPDK: %v", err))
	}

	handle, err := dpdk.NewDPDKHandle(0, "")
	if err != nil {
		panic(fmt.Sprintf("Failed to create DPDK handle: %v", err))
	}
	defer handle.Close()
	handle.PrintInfo()
	if !handle.IsPortUp() {
		panic("Port is not up")
	}

	fmt.Println("Starting packet capture...")
	for {
		packet, _, err := handle.ReadPacketData()
		if err != nil {
			fmt.Printf("Error reading packet: %v\n", err)
			continue
		}
		if packet == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		packet_ := gopacket_dpdk.NewPacket(packet, layers.LayerTypeEthernet, gopacket_dpdk.Default)
		ipLayer := packet_.Layer(layers.LayerTypeIPv4)
		if ipLayer != nil {
			ip, ok := ipLayer.(*layers.IPv4)
			if !ok {
				fmt.Println("failed to parse IPv4 layer")
				return
			}
			fmt.Printf("Source IP: %s, Destination IP: %s\n", ip.SrcIP, ip.DstIP)
		}
	}
}
