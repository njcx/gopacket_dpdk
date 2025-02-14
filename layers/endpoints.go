// Copyright 2012 Google, Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

package layers

import (
	"encoding/binary"
	"github.com/njcx/gopacket_dpdk"
	"net"
	"strconv"
)

var (
	// We use two different endpoint types for IPv4 vs IPv6 addresses, so that
	// ordering with endpointA.LessThan(endpointB) sanely groups all IPv4
	// addresses and all IPv6 addresses, such that IPv6 > IPv4 for all addresses.
	EndpointIPv4 = gopacket_dpdk.RegisterEndpointType(1, gopacket_dpdk.EndpointTypeMetadata{"IPv4", func(b []byte) string {
		return net.IP(b).String()
	}})
	EndpointIPv6 = gopacket_dpdk.RegisterEndpointType(2, gopacket_dpdk.EndpointTypeMetadata{"IPv6", func(b []byte) string {
		return net.IP(b).String()
	}})

	EndpointMAC = gopacket_dpdk.RegisterEndpointType(3, gopacket_dpdk.EndpointTypeMetadata{"MAC", func(b []byte) string {
		return net.HardwareAddr(b).String()
	}})
	EndpointTCPPort = gopacket_dpdk.RegisterEndpointType(4, gopacket_dpdk.EndpointTypeMetadata{"TCP", func(b []byte) string {
		return strconv.Itoa(int(binary.BigEndian.Uint16(b)))
	}})
	EndpointUDPPort = gopacket_dpdk.RegisterEndpointType(5, gopacket_dpdk.EndpointTypeMetadata{"UDP", func(b []byte) string {
		return strconv.Itoa(int(binary.BigEndian.Uint16(b)))
	}})
	EndpointSCTPPort = gopacket_dpdk.RegisterEndpointType(6, gopacket_dpdk.EndpointTypeMetadata{"SCTP", func(b []byte) string {
		return strconv.Itoa(int(binary.BigEndian.Uint16(b)))
	}})
	EndpointRUDPPort = gopacket_dpdk.RegisterEndpointType(7, gopacket_dpdk.EndpointTypeMetadata{"RUDP", func(b []byte) string {
		return strconv.Itoa(int(b[0]))
	}})
	EndpointUDPLitePort = gopacket_dpdk.RegisterEndpointType(8, gopacket_dpdk.EndpointTypeMetadata{"UDPLite", func(b []byte) string {
		return strconv.Itoa(int(binary.BigEndian.Uint16(b)))
	}})
	EndpointPPP = gopacket_dpdk.RegisterEndpointType(9, gopacket_dpdk.EndpointTypeMetadata{"PPP", func([]byte) string {
		return "point"
	}})
)

// NewIPEndpoint creates a new IP (v4 or v6) endpoint from a net.IP address.
// It returns gopacket_dpdk.InvalidEndpoint if the IP address is invalid.
func NewIPEndpoint(a net.IP) gopacket_dpdk.Endpoint {
	switch len(a) {
	case 4:
		return gopacket_dpdk.NewEndpoint(EndpointIPv4, []byte(a))
	case 16:
		return gopacket_dpdk.NewEndpoint(EndpointIPv6, []byte(a))
	}
	return gopacket_dpdk.InvalidEndpoint
}

// NewMACEndpoint returns a new MAC address endpoint.
func NewMACEndpoint(a net.HardwareAddr) gopacket_dpdk.Endpoint {
	return gopacket_dpdk.NewEndpoint(EndpointMAC, []byte(a))
}
func newPortEndpoint(t gopacket_dpdk.EndpointType, p uint16) gopacket_dpdk.Endpoint {
	return gopacket_dpdk.NewEndpoint(t, []byte{byte(p >> 8), byte(p)})
}

// NewTCPPortEndpoint returns an endpoint based on a TCP port.
func NewTCPPortEndpoint(p TCPPort) gopacket_dpdk.Endpoint {
	return newPortEndpoint(EndpointTCPPort, uint16(p))
}

// NewUDPPortEndpoint returns an endpoint based on a UDP port.
func NewUDPPortEndpoint(p UDPPort) gopacket_dpdk.Endpoint {
	return newPortEndpoint(EndpointUDPPort, uint16(p))
}

// NewSCTPPortEndpoint returns an endpoint based on a SCTP port.
func NewSCTPPortEndpoint(p SCTPPort) gopacket_dpdk.Endpoint {
	return newPortEndpoint(EndpointSCTPPort, uint16(p))
}

// NewRUDPPortEndpoint returns an endpoint based on a RUDP port.
func NewRUDPPortEndpoint(p RUDPPort) gopacket_dpdk.Endpoint {
	return gopacket_dpdk.NewEndpoint(EndpointRUDPPort, []byte{byte(p)})
}

// NewUDPLitePortEndpoint returns an endpoint based on a UDPLite port.
func NewUDPLitePortEndpoint(p UDPLitePort) gopacket_dpdk.Endpoint {
	return newPortEndpoint(EndpointUDPLitePort, uint16(p))
}
