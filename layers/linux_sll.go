// Copyright 2012 Google, Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

package layers

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"github.com/njcx/gopacket_dpdk"
)

type LinuxSLLPacketType uint16

const (
	LinuxSLLPacketTypeHost      LinuxSLLPacketType = 0 // To us
	LinuxSLLPacketTypeBroadcast LinuxSLLPacketType = 1 // To all
	LinuxSLLPacketTypeMulticast LinuxSLLPacketType = 2 // To group
	LinuxSLLPacketTypeOtherhost LinuxSLLPacketType = 3 // To someone else
	LinuxSLLPacketTypeOutgoing  LinuxSLLPacketType = 4 // Outgoing of any type
	// These ones are invisible by user level
	LinuxSLLPacketTypeLoopback  LinuxSLLPacketType = 5 // MC/BRD frame looped back
	LinuxSLLPacketTypeFastroute LinuxSLLPacketType = 6 // Fastrouted frame
)

func (l LinuxSLLPacketType) String() string {
	switch l {
	case LinuxSLLPacketTypeHost:
		return "host"
	case LinuxSLLPacketTypeBroadcast:
		return "broadcast"
	case LinuxSLLPacketTypeMulticast:
		return "multicast"
	case LinuxSLLPacketTypeOtherhost:
		return "otherhost"
	case LinuxSLLPacketTypeOutgoing:
		return "outgoing"
	case LinuxSLLPacketTypeLoopback:
		return "loopback"
	case LinuxSLLPacketTypeFastroute:
		return "fastroute"
	}
	return fmt.Sprintf("Unknown(%d)", int(l))
}

type LinuxSLL struct {
	BaseLayer
	PacketType   LinuxSLLPacketType
	AddrLen      uint16
	Addr         net.HardwareAddr
	EthernetType EthernetType
}

// LayerType returns LayerTypeLinuxSLL.
func (sll *LinuxSLL) LayerType() gopacket_dpdk.LayerType { return LayerTypeLinuxSLL }

func (sll *LinuxSLL) CanDecode() gopacket_dpdk.LayerClass {
	return LayerTypeLinuxSLL
}

func (sll *LinuxSLL) LinkFlow() gopacket_dpdk.Flow {
	return gopacket_dpdk.NewFlow(EndpointMAC, sll.Addr, nil)
}

func (sll *LinuxSLL) NextLayerType() gopacket_dpdk.LayerType {
	return sll.EthernetType.LayerType()
}

func (sll *LinuxSLL) DecodeFromBytes(data []byte, df gopacket_dpdk.DecodeFeedback) error {
	if len(data) < 16 {
		return errors.New("Linux SLL packet too small")
	}
	sll.PacketType = LinuxSLLPacketType(binary.BigEndian.Uint16(data[0:2]))
	sll.AddrLen = binary.BigEndian.Uint16(data[4:6])

	sll.Addr = net.HardwareAddr(data[6 : sll.AddrLen+6])
	sll.EthernetType = EthernetType(binary.BigEndian.Uint16(data[14:16]))
	sll.BaseLayer = BaseLayer{data[:16], data[16:]}

	return nil
}

func decodeLinuxSLL(data []byte, p gopacket_dpdk.PacketBuilder) error {
	sll := &LinuxSLL{}
	if err := sll.DecodeFromBytes(data, p); err != nil {
		return err
	}
	p.AddLayer(sll)
	p.SetLinkLayer(sll)
	return p.NextDecoder(sll.EthernetType)
}
