package main

import (
	"fmt"
	"github.com/njcx/gopacket_dpdk"
	"github.com/njcx/gopacket_dpdk/dpdk"
	"time"
)

func main() {
	// Initialize DPDK
	if err := dpdk.InitDPDK(); err != nil {
		panic(fmt.Sprintf("Failed to initialize DPDK: %v", err))
	}
	// Create a new DPDK handle for the first port
	handle, err := dpdk.NewDPDKHandle(0)
	if err != nil {
		panic(fmt.Sprintf("Failed to create DPDK handle: %v", err))
	}
	defer handle.Close()

	// Print port information
	handle.PrintInfo()

	// Check if port is up
	if !handle.IsPortUp() {
		panic("Port is not up")
	}

	// Read packets in a loop
	fmt.Println("Starting packet capture...")
	for i := 0; i < 100; i++ {
		packet, err := handle.ReadPacket()
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
