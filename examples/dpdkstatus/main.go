package main

import (
	"fmt"
	"github.com/njcx/gopacket_dpdk/dpdk"
	"time"
)

func main() {
	// Initialize DPDK
	if err := dpdk.InitDPDK(); err != nil {
		panic(fmt.Sprintf("Failed to initialize DPDK: %v", err))
	}

	// Create a new DPDK handle for the first port
	handle, err := dpdk.NewDPDKHandle(0, "")
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

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			fmt.Println("\n=== Statistics Update ===")

			if err := handle.PrintBandwidth(); err != nil {
				fmt.Printf("Error printing bandwidth: %v\n", err)
			}

			if err := handle.PrintPacketLoss(); err != nil {
				fmt.Printf("Error printing packet loss: %v\n", err)
			}
		}
	}()

}
