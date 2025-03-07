// Copyright 2012 Google, Inc. All rights reserved.
// Copyright 2009-2011 Andreas Krennmair. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

package layers

import (
	"fmt"
	"github.com/njcx/gopacket_dpdk"
)

// Checksum computation for TCP/UDP.
type tcpipchecksum struct {
	pseudoheader tcpipPseudoHeader
}

type tcpipPseudoHeader interface {
	pseudoheaderChecksum() uint32
}

func (ip *IPv4) pseudoheaderChecksum() (csum uint32) {
	csum += (uint32(ip.SrcIP[0]) + uint32(ip.SrcIP[2])) << 8
	csum += uint32(ip.SrcIP[1]) + uint32(ip.SrcIP[3])
	csum += (uint32(ip.DstIP[0]) + uint32(ip.DstIP[2])) << 8
	csum += uint32(ip.DstIP[1]) + uint32(ip.DstIP[3])
	csum += uint32(ip.Protocol)
	return
}

func (ip *IPv6) pseudoheaderChecksum() (csum uint32) {
	for i := 0; i < 16; i += 2 {
		csum += uint32(ip.SrcIP[i]) << 8
		csum += uint32(ip.SrcIP[i+1])
		csum += uint32(ip.DstIP[i]) << 8
		csum += uint32(ip.DstIP[i+1])
	}
	csum += uint32(ip.NextHeader)
	return
}

// Calculate the TCP/IP checksum defined in rfc1071.  The passed-in csum is any
// initial checksum data that's already been computed.
func tcpipChecksum(data []byte, csum uint32) uint16 {
	// to handle odd lengths, we loop to length - 1, incrementing by 2, then
	// handle the last byte specifically by checking against the original
	// length.
	length := len(data) - 1
	for i := 0; i < length; i += 2 {
		// For our test packet, doing this manually is about 25% faster
		// (740 ns vs. 1000ns) than doing it by calling binary.BigEndian.Uint16.
		csum += uint32(data[i]) << 8
		csum += uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		csum += uint32(data[length]) << 8
	}
	for csum > 0xffff {
		csum = (csum >> 16) + (csum & 0xffff)
	}
	return ^uint16(csum + (csum >> 16))
}

// computeChecksum computes a TCP or UDP checksum.  headerAndPayload is the
// serialized TCP or UDP header plus its payload, with the checksum zero'd
// out.
func (c *tcpipchecksum) computeChecksum(headerAndPayload []byte) (uint16, error) {
	if c.pseudoheader == nil {
		return 0, fmt.Errorf("TCP/IP layer 4 checksum cannot be computed without network layer... call SetNetworkLayerForChecksum to set which layer to use")
	}
	length := uint32(len(headerAndPayload))
	csum := c.pseudoheader.pseudoheaderChecksum()
	csum += length & 0xffff
	csum += length >> 16
	return tcpipChecksum(headerAndPayload, csum), nil
}

// SetNetworkLayerForChecksum tells this layer which network layer is wrapping it.
// This is needed for computing the checksum when serializing, since TCP/IP transport
// layer checksums depends on fields in the IPv4 or IPv6 layer that contains it.
// The passed in layer must be an *IPv4 or *IPv6.
func (i *tcpipchecksum) SetNetworkLayerForChecksum(l gopacket_dpdk.NetworkLayer) error {
	switch v := l.(type) {
	case *IPv4:
		i.pseudoheader = v
	case *IPv6:
		i.pseudoheader = v
	default:
		return fmt.Errorf("cannot use layer type %v for tcp checksum network layer", l.LayerType())
	}
	return nil
}
