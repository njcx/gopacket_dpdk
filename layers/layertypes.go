// Copyright 2012 Google, gopacket_dpdk.LayerTypeMetadata{Inc. All rights reserved}.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

package layers

import (
	"github.com/njcx/gopacket_dpdk"
)

var (
	LayerTypeARP                         = gopacket_dpdk.RegisterLayerType(10, gopacket_dpdk.LayerTypeMetadata{"ARP", gopacket_dpdk.DecodeFunc(decodeARP)})
	LayerTypeCiscoDiscovery              = gopacket_dpdk.RegisterLayerType(11, gopacket_dpdk.LayerTypeMetadata{"CiscoDiscovery", gopacket_dpdk.DecodeFunc(decodeCiscoDiscovery)})
	LayerTypeEthernetCTP                 = gopacket_dpdk.RegisterLayerType(12, gopacket_dpdk.LayerTypeMetadata{"EthernetCTP", gopacket_dpdk.DecodeFunc(decodeEthernetCTP)})
	LayerTypeEthernetCTPForwardData      = gopacket_dpdk.RegisterLayerType(13, gopacket_dpdk.LayerTypeMetadata{"EthernetCTPForwardData", nil})
	LayerTypeEthernetCTPReply            = gopacket_dpdk.RegisterLayerType(14, gopacket_dpdk.LayerTypeMetadata{"EthernetCTPReply", nil})
	LayerTypeDot1Q                       = gopacket_dpdk.RegisterLayerType(15, gopacket_dpdk.LayerTypeMetadata{"Dot1Q", gopacket_dpdk.DecodeFunc(decodeDot1Q)})
	LayerTypeEtherIP                     = gopacket_dpdk.RegisterLayerType(16, gopacket_dpdk.LayerTypeMetadata{"EtherIP", gopacket_dpdk.DecodeFunc(decodeEtherIP)})
	LayerTypeEthernet                    = gopacket_dpdk.RegisterLayerType(17, gopacket_dpdk.LayerTypeMetadata{"Ethernet", gopacket_dpdk.DecodeFunc(decodeEthernet)})
	LayerTypeGRE                         = gopacket_dpdk.RegisterLayerType(18, gopacket_dpdk.LayerTypeMetadata{"GRE", gopacket_dpdk.DecodeFunc(decodeGRE)})
	LayerTypeICMPv4                      = gopacket_dpdk.RegisterLayerType(19, gopacket_dpdk.LayerTypeMetadata{"ICMPv4", gopacket_dpdk.DecodeFunc(decodeICMPv4)})
	LayerTypeIPv4                        = gopacket_dpdk.RegisterLayerType(20, gopacket_dpdk.LayerTypeMetadata{"IPv4", gopacket_dpdk.DecodeFunc(decodeIPv4)})
	LayerTypeIPv6                        = gopacket_dpdk.RegisterLayerType(21, gopacket_dpdk.LayerTypeMetadata{"IPv6", gopacket_dpdk.DecodeFunc(decodeIPv6)})
	LayerTypeLLC                         = gopacket_dpdk.RegisterLayerType(22, gopacket_dpdk.LayerTypeMetadata{"LLC", gopacket_dpdk.DecodeFunc(decodeLLC)})
	LayerTypeSNAP                        = gopacket_dpdk.RegisterLayerType(23, gopacket_dpdk.LayerTypeMetadata{"SNAP", gopacket_dpdk.DecodeFunc(decodeSNAP)})
	LayerTypeMPLS                        = gopacket_dpdk.RegisterLayerType(24, gopacket_dpdk.LayerTypeMetadata{"MPLS", gopacket_dpdk.DecodeFunc(decodeMPLS)})
	LayerTypePPP                         = gopacket_dpdk.RegisterLayerType(25, gopacket_dpdk.LayerTypeMetadata{"PPP", gopacket_dpdk.DecodeFunc(decodePPP)})
	LayerTypePPPoE                       = gopacket_dpdk.RegisterLayerType(26, gopacket_dpdk.LayerTypeMetadata{"PPPoE", gopacket_dpdk.DecodeFunc(decodePPPoE)})
	LayerTypeRUDP                        = gopacket_dpdk.RegisterLayerType(27, gopacket_dpdk.LayerTypeMetadata{"RUDP", gopacket_dpdk.DecodeFunc(decodeRUDP)})
	LayerTypeSCTP                        = gopacket_dpdk.RegisterLayerType(28, gopacket_dpdk.LayerTypeMetadata{"SCTP", gopacket_dpdk.DecodeFunc(decodeSCTP)})
	LayerTypeSCTPUnknownChunkType        = gopacket_dpdk.RegisterLayerType(29, gopacket_dpdk.LayerTypeMetadata{"SCTPUnknownChunkType", nil})
	LayerTypeSCTPData                    = gopacket_dpdk.RegisterLayerType(30, gopacket_dpdk.LayerTypeMetadata{"SCTPData", nil})
	LayerTypeSCTPInit                    = gopacket_dpdk.RegisterLayerType(31, gopacket_dpdk.LayerTypeMetadata{"SCTPInit", nil})
	LayerTypeSCTPSack                    = gopacket_dpdk.RegisterLayerType(32, gopacket_dpdk.LayerTypeMetadata{"SCTPSack", nil})
	LayerTypeSCTPHeartbeat               = gopacket_dpdk.RegisterLayerType(33, gopacket_dpdk.LayerTypeMetadata{"SCTPHeartbeat", nil})
	LayerTypeSCTPError                   = gopacket_dpdk.RegisterLayerType(34, gopacket_dpdk.LayerTypeMetadata{"SCTPError", nil})
	LayerTypeSCTPShutdown                = gopacket_dpdk.RegisterLayerType(35, gopacket_dpdk.LayerTypeMetadata{"SCTPShutdown", nil})
	LayerTypeSCTPShutdownAck             = gopacket_dpdk.RegisterLayerType(36, gopacket_dpdk.LayerTypeMetadata{"SCTPShutdownAck", nil})
	LayerTypeSCTPCookieEcho              = gopacket_dpdk.RegisterLayerType(37, gopacket_dpdk.LayerTypeMetadata{"SCTPCookieEcho", nil})
	LayerTypeSCTPEmptyLayer              = gopacket_dpdk.RegisterLayerType(38, gopacket_dpdk.LayerTypeMetadata{"SCTPEmptyLayer", nil})
	LayerTypeSCTPInitAck                 = gopacket_dpdk.RegisterLayerType(39, gopacket_dpdk.LayerTypeMetadata{"SCTPInitAck", nil})
	LayerTypeSCTPHeartbeatAck            = gopacket_dpdk.RegisterLayerType(40, gopacket_dpdk.LayerTypeMetadata{"SCTPHeartbeatAck", nil})
	LayerTypeSCTPAbort                   = gopacket_dpdk.RegisterLayerType(41, gopacket_dpdk.LayerTypeMetadata{"SCTPAbort", nil})
	LayerTypeSCTPShutdownComplete        = gopacket_dpdk.RegisterLayerType(42, gopacket_dpdk.LayerTypeMetadata{"SCTPShutdownComplete", nil})
	LayerTypeSCTPCookieAck               = gopacket_dpdk.RegisterLayerType(43, gopacket_dpdk.LayerTypeMetadata{"SCTPCookieAck", nil})
	LayerTypeTCP                         = gopacket_dpdk.RegisterLayerType(44, gopacket_dpdk.LayerTypeMetadata{"TCP", gopacket_dpdk.DecodeFunc(decodeTCP)})
	LayerTypeUDP                         = gopacket_dpdk.RegisterLayerType(45, gopacket_dpdk.LayerTypeMetadata{"UDP", gopacket_dpdk.DecodeFunc(decodeUDP)})
	LayerTypeIPv6HopByHop                = gopacket_dpdk.RegisterLayerType(46, gopacket_dpdk.LayerTypeMetadata{"IPv6HopByHop", gopacket_dpdk.DecodeFunc(decodeIPv6HopByHop)})
	LayerTypeIPv6Routing                 = gopacket_dpdk.RegisterLayerType(47, gopacket_dpdk.LayerTypeMetadata{"IPv6Routing", gopacket_dpdk.DecodeFunc(decodeIPv6Routing)})
	LayerTypeIPv6Fragment                = gopacket_dpdk.RegisterLayerType(48, gopacket_dpdk.LayerTypeMetadata{"IPv6Fragment", gopacket_dpdk.DecodeFunc(decodeIPv6Fragment)})
	LayerTypeIPv6Destination             = gopacket_dpdk.RegisterLayerType(49, gopacket_dpdk.LayerTypeMetadata{"IPv6Destination", gopacket_dpdk.DecodeFunc(decodeIPv6Destination)})
	LayerTypeIPSecAH                     = gopacket_dpdk.RegisterLayerType(50, gopacket_dpdk.LayerTypeMetadata{"IPSecAH", gopacket_dpdk.DecodeFunc(decodeIPSecAH)})
	LayerTypeIPSecESP                    = gopacket_dpdk.RegisterLayerType(51, gopacket_dpdk.LayerTypeMetadata{"IPSecESP", gopacket_dpdk.DecodeFunc(decodeIPSecESP)})
	LayerTypeUDPLite                     = gopacket_dpdk.RegisterLayerType(52, gopacket_dpdk.LayerTypeMetadata{"UDPLite", gopacket_dpdk.DecodeFunc(decodeUDPLite)})
	LayerTypeFDDI                        = gopacket_dpdk.RegisterLayerType(53, gopacket_dpdk.LayerTypeMetadata{"FDDI", gopacket_dpdk.DecodeFunc(decodeFDDI)})
	LayerTypeLoopback                    = gopacket_dpdk.RegisterLayerType(54, gopacket_dpdk.LayerTypeMetadata{"Loopback", gopacket_dpdk.DecodeFunc(decodeLoopback)})
	LayerTypeEAP                         = gopacket_dpdk.RegisterLayerType(55, gopacket_dpdk.LayerTypeMetadata{"EAP", gopacket_dpdk.DecodeFunc(decodeEAP)})
	LayerTypeEAPOL                       = gopacket_dpdk.RegisterLayerType(56, gopacket_dpdk.LayerTypeMetadata{"EAPOL", gopacket_dpdk.DecodeFunc(decodeEAPOL)})
	LayerTypeICMPv6                      = gopacket_dpdk.RegisterLayerType(57, gopacket_dpdk.LayerTypeMetadata{"ICMPv6", gopacket_dpdk.DecodeFunc(decodeICMPv6)})
	LayerTypeLinkLayerDiscovery          = gopacket_dpdk.RegisterLayerType(58, gopacket_dpdk.LayerTypeMetadata{"LinkLayerDiscovery", gopacket_dpdk.DecodeFunc(decodeLinkLayerDiscovery)})
	LayerTypeCiscoDiscoveryInfo          = gopacket_dpdk.RegisterLayerType(59, gopacket_dpdk.LayerTypeMetadata{"CiscoDiscoveryInfo", gopacket_dpdk.DecodeFunc(decodeCiscoDiscoveryInfo)})
	LayerTypeLinkLayerDiscoveryInfo      = gopacket_dpdk.RegisterLayerType(60, gopacket_dpdk.LayerTypeMetadata{"LinkLayerDiscoveryInfo", nil})
	LayerTypeNortelDiscovery             = gopacket_dpdk.RegisterLayerType(61, gopacket_dpdk.LayerTypeMetadata{"NortelDiscovery", gopacket_dpdk.DecodeFunc(decodeNortelDiscovery)})
	LayerTypeIGMP                        = gopacket_dpdk.RegisterLayerType(62, gopacket_dpdk.LayerTypeMetadata{"IGMP", gopacket_dpdk.DecodeFunc(decodeIGMP)})
	LayerTypePFLog                       = gopacket_dpdk.RegisterLayerType(63, gopacket_dpdk.LayerTypeMetadata{"PFLog", gopacket_dpdk.DecodeFunc(decodePFLog)})
	LayerTypeRadioTap                    = gopacket_dpdk.RegisterLayerType(64, gopacket_dpdk.LayerTypeMetadata{"RadioTap", gopacket_dpdk.DecodeFunc(decodeRadioTap)})
	LayerTypeDot11                       = gopacket_dpdk.RegisterLayerType(65, gopacket_dpdk.LayerTypeMetadata{"Dot11", gopacket_dpdk.DecodeFunc(decodeDot11)})
	LayerTypeDot11Ctrl                   = gopacket_dpdk.RegisterLayerType(66, gopacket_dpdk.LayerTypeMetadata{"Dot11Ctrl", gopacket_dpdk.DecodeFunc(decodeDot11Ctrl)})
	LayerTypeDot11Data                   = gopacket_dpdk.RegisterLayerType(67, gopacket_dpdk.LayerTypeMetadata{"Dot11Data", gopacket_dpdk.DecodeFunc(decodeDot11Data)})
	LayerTypeDot11DataCFAck              = gopacket_dpdk.RegisterLayerType(68, gopacket_dpdk.LayerTypeMetadata{"Dot11DataCFAck", gopacket_dpdk.DecodeFunc(decodeDot11DataCFAck)})
	LayerTypeDot11DataCFPoll             = gopacket_dpdk.RegisterLayerType(69, gopacket_dpdk.LayerTypeMetadata{"Dot11DataCFPoll", gopacket_dpdk.DecodeFunc(decodeDot11DataCFPoll)})
	LayerTypeDot11DataCFAckPoll          = gopacket_dpdk.RegisterLayerType(70, gopacket_dpdk.LayerTypeMetadata{"Dot11DataCFAckPoll", gopacket_dpdk.DecodeFunc(decodeDot11DataCFAckPoll)})
	LayerTypeDot11DataNull               = gopacket_dpdk.RegisterLayerType(71, gopacket_dpdk.LayerTypeMetadata{"Dot11DataNull", gopacket_dpdk.DecodeFunc(decodeDot11DataNull)})
	LayerTypeDot11DataCFAckNoData        = gopacket_dpdk.RegisterLayerType(72, gopacket_dpdk.LayerTypeMetadata{"Dot11DataCFAck", gopacket_dpdk.DecodeFunc(decodeDot11DataCFAck)})
	LayerTypeDot11DataCFPollNoData       = gopacket_dpdk.RegisterLayerType(73, gopacket_dpdk.LayerTypeMetadata{"Dot11DataCFPoll", gopacket_dpdk.DecodeFunc(decodeDot11DataCFPoll)})
	LayerTypeDot11DataCFAckPollNoData    = gopacket_dpdk.RegisterLayerType(74, gopacket_dpdk.LayerTypeMetadata{"Dot11DataCFAckPoll", gopacket_dpdk.DecodeFunc(decodeDot11DataCFAckPoll)})
	LayerTypeDot11DataQOSData            = gopacket_dpdk.RegisterLayerType(75, gopacket_dpdk.LayerTypeMetadata{"Dot11DataQOSData", gopacket_dpdk.DecodeFunc(decodeDot11DataQOSData)})
	LayerTypeDot11DataQOSDataCFAck       = gopacket_dpdk.RegisterLayerType(76, gopacket_dpdk.LayerTypeMetadata{"Dot11DataQOSDataCFAck", gopacket_dpdk.DecodeFunc(decodeDot11DataQOSDataCFAck)})
	LayerTypeDot11DataQOSDataCFPoll      = gopacket_dpdk.RegisterLayerType(77, gopacket_dpdk.LayerTypeMetadata{"Dot11DataQOSDataCFPoll", gopacket_dpdk.DecodeFunc(decodeDot11DataQOSDataCFPoll)})
	LayerTypeDot11DataQOSDataCFAckPoll   = gopacket_dpdk.RegisterLayerType(78, gopacket_dpdk.LayerTypeMetadata{"Dot11DataQOSDataCFAckPoll", gopacket_dpdk.DecodeFunc(decodeDot11DataQOSDataCFAckPoll)})
	LayerTypeDot11DataQOSNull            = gopacket_dpdk.RegisterLayerType(79, gopacket_dpdk.LayerTypeMetadata{"Dot11DataQOSNull", gopacket_dpdk.DecodeFunc(decodeDot11DataQOSNull)})
	LayerTypeDot11DataQOSCFPollNoData    = gopacket_dpdk.RegisterLayerType(80, gopacket_dpdk.LayerTypeMetadata{"Dot11DataQOSCFPoll", gopacket_dpdk.DecodeFunc(decodeDot11DataQOSCFPollNoData)})
	LayerTypeDot11DataQOSCFAckPollNoData = gopacket_dpdk.RegisterLayerType(81, gopacket_dpdk.LayerTypeMetadata{"Dot11DataQOSCFAckPoll", gopacket_dpdk.DecodeFunc(decodeDot11DataQOSCFAckPollNoData)})
	LayerTypeDot11InformationElement     = gopacket_dpdk.RegisterLayerType(82, gopacket_dpdk.LayerTypeMetadata{"Dot11InformationElement", gopacket_dpdk.DecodeFunc(decodeDot11InformationElement)})
	LayerTypeDot11CtrlCTS                = gopacket_dpdk.RegisterLayerType(83, gopacket_dpdk.LayerTypeMetadata{"Dot11CtrlCTS", gopacket_dpdk.DecodeFunc(decodeDot11CtrlCTS)})
	LayerTypeDot11CtrlRTS                = gopacket_dpdk.RegisterLayerType(84, gopacket_dpdk.LayerTypeMetadata{"Dot11CtrlRTS", gopacket_dpdk.DecodeFunc(decodeDot11CtrlRTS)})
	LayerTypeDot11CtrlBlockAckReq        = gopacket_dpdk.RegisterLayerType(85, gopacket_dpdk.LayerTypeMetadata{"Dot11CtrlBlockAckReq", gopacket_dpdk.DecodeFunc(decodeDot11CtrlBlockAckReq)})
	LayerTypeDot11CtrlBlockAck           = gopacket_dpdk.RegisterLayerType(86, gopacket_dpdk.LayerTypeMetadata{"Dot11CtrlBlockAck", gopacket_dpdk.DecodeFunc(decodeDot11CtrlBlockAck)})
	LayerTypeDot11CtrlPowersavePoll      = gopacket_dpdk.RegisterLayerType(87, gopacket_dpdk.LayerTypeMetadata{"Dot11CtrlPowersavePoll", gopacket_dpdk.DecodeFunc(decodeDot11CtrlPowersavePoll)})
	LayerTypeDot11CtrlAck                = gopacket_dpdk.RegisterLayerType(88, gopacket_dpdk.LayerTypeMetadata{"Dot11CtrlAck", gopacket_dpdk.DecodeFunc(decodeDot11CtrlAck)})
	LayerTypeDot11CtrlCFEnd              = gopacket_dpdk.RegisterLayerType(89, gopacket_dpdk.LayerTypeMetadata{"Dot11CtrlCFEnd", gopacket_dpdk.DecodeFunc(decodeDot11CtrlCFEnd)})
	LayerTypeDot11CtrlCFEndAck           = gopacket_dpdk.RegisterLayerType(90, gopacket_dpdk.LayerTypeMetadata{"Dot11CtrlCFEndAck", gopacket_dpdk.DecodeFunc(decodeDot11CtrlCFEndAck)})
	LayerTypeDot11MgmtAssociationReq     = gopacket_dpdk.RegisterLayerType(91, gopacket_dpdk.LayerTypeMetadata{"Dot11MgmtAssociationReq", gopacket_dpdk.DecodeFunc(decodeDot11MgmtAssociationReq)})
	LayerTypeDot11MgmtAssociationResp    = gopacket_dpdk.RegisterLayerType(92, gopacket_dpdk.LayerTypeMetadata{"Dot11MgmtAssociationResp", gopacket_dpdk.DecodeFunc(decodeDot11MgmtAssociationResp)})
	LayerTypeDot11MgmtReassociationReq   = gopacket_dpdk.RegisterLayerType(93, gopacket_dpdk.LayerTypeMetadata{"Dot11MgmtReassociationReq", gopacket_dpdk.DecodeFunc(decodeDot11MgmtReassociationReq)})
	LayerTypeDot11MgmtReassociationResp  = gopacket_dpdk.RegisterLayerType(94, gopacket_dpdk.LayerTypeMetadata{"Dot11MgmtReassociationResp", gopacket_dpdk.DecodeFunc(decodeDot11MgmtReassociationResp)})
	LayerTypeDot11MgmtProbeReq           = gopacket_dpdk.RegisterLayerType(95, gopacket_dpdk.LayerTypeMetadata{"Dot11MgmtProbeReq", gopacket_dpdk.DecodeFunc(decodeDot11MgmtProbeReq)})
	LayerTypeDot11MgmtProbeResp          = gopacket_dpdk.RegisterLayerType(96, gopacket_dpdk.LayerTypeMetadata{"Dot11MgmtProbeResp", gopacket_dpdk.DecodeFunc(decodeDot11MgmtProbeResp)})
	LayerTypeDot11MgmtMeasurementPilot   = gopacket_dpdk.RegisterLayerType(97, gopacket_dpdk.LayerTypeMetadata{"Dot11MgmtMeasurementPilot", gopacket_dpdk.DecodeFunc(decodeDot11MgmtMeasurementPilot)})
	LayerTypeDot11MgmtBeacon             = gopacket_dpdk.RegisterLayerType(98, gopacket_dpdk.LayerTypeMetadata{"Dot11MgmtBeacon", gopacket_dpdk.DecodeFunc(decodeDot11MgmtBeacon)})
	LayerTypeDot11MgmtATIM               = gopacket_dpdk.RegisterLayerType(99, gopacket_dpdk.LayerTypeMetadata{"Dot11MgmtATIM", gopacket_dpdk.DecodeFunc(decodeDot11MgmtATIM)})
	LayerTypeDot11MgmtDisassociation     = gopacket_dpdk.RegisterLayerType(100, gopacket_dpdk.LayerTypeMetadata{"Dot11MgmtDisassociation", gopacket_dpdk.DecodeFunc(decodeDot11MgmtDisassociation)})
	LayerTypeDot11MgmtAuthentication     = gopacket_dpdk.RegisterLayerType(101, gopacket_dpdk.LayerTypeMetadata{"Dot11MgmtAuthentication", gopacket_dpdk.DecodeFunc(decodeDot11MgmtAuthentication)})
	LayerTypeDot11MgmtDeauthentication   = gopacket_dpdk.RegisterLayerType(102, gopacket_dpdk.LayerTypeMetadata{"Dot11MgmtDeauthentication", gopacket_dpdk.DecodeFunc(decodeDot11MgmtDeauthentication)})
	LayerTypeDot11MgmtAction             = gopacket_dpdk.RegisterLayerType(103, gopacket_dpdk.LayerTypeMetadata{"Dot11MgmtAction", gopacket_dpdk.DecodeFunc(decodeDot11MgmtAction)})
	LayerTypeDot11MgmtActionNoAck        = gopacket_dpdk.RegisterLayerType(104, gopacket_dpdk.LayerTypeMetadata{"Dot11MgmtActionNoAck", gopacket_dpdk.DecodeFunc(decodeDot11MgmtActionNoAck)})
	LayerTypeDot11MgmtArubaWLAN          = gopacket_dpdk.RegisterLayerType(105, gopacket_dpdk.LayerTypeMetadata{"Dot11MgmtArubaWLAN", gopacket_dpdk.DecodeFunc(decodeDot11MgmtArubaWLAN)})
	LayerTypeDot11WEP                    = gopacket_dpdk.RegisterLayerType(106, gopacket_dpdk.LayerTypeMetadata{"Dot11WEP", gopacket_dpdk.DecodeFunc(decodeDot11WEP)})
	LayerTypeDNS                         = gopacket_dpdk.RegisterLayerType(107, gopacket_dpdk.LayerTypeMetadata{"DNS", gopacket_dpdk.DecodeFunc(decodeDNS)})
	LayerTypeUSB                         = gopacket_dpdk.RegisterLayerType(108, gopacket_dpdk.LayerTypeMetadata{"USB", gopacket_dpdk.DecodeFunc(decodeUSB)})
	LayerTypeUSBRequestBlockSetup        = gopacket_dpdk.RegisterLayerType(109, gopacket_dpdk.LayerTypeMetadata{"USBRequestBlockSetup", gopacket_dpdk.DecodeFunc(decodeUSBRequestBlockSetup)})
	LayerTypeUSBControl                  = gopacket_dpdk.RegisterLayerType(110, gopacket_dpdk.LayerTypeMetadata{"USBControl", gopacket_dpdk.DecodeFunc(decodeUSBControl)})
	LayerTypeUSBInterrupt                = gopacket_dpdk.RegisterLayerType(111, gopacket_dpdk.LayerTypeMetadata{"USBInterrupt", gopacket_dpdk.DecodeFunc(decodeUSBInterrupt)})
	LayerTypeUSBBulk                     = gopacket_dpdk.RegisterLayerType(112, gopacket_dpdk.LayerTypeMetadata{"USBBulk", gopacket_dpdk.DecodeFunc(decodeUSBBulk)})
	LayerTypeLinuxSLL                    = gopacket_dpdk.RegisterLayerType(113, gopacket_dpdk.LayerTypeMetadata{"Linux SLL", gopacket_dpdk.DecodeFunc(decodeLinuxSLL)})
)

var (
	// LayerClassIPNetwork contains TCP/IP network layer types.
	LayerClassIPNetwork = gopacket_dpdk.NewLayerClass([]gopacket_dpdk.LayerType{
		LayerTypeIPv4,
		LayerTypeIPv6,
	})
	// LayerClassIPTransport contains TCP/IP transport layer types.
	LayerClassIPTransport = gopacket_dpdk.NewLayerClass([]gopacket_dpdk.LayerType{
		LayerTypeTCP,
		LayerTypeUDP,
		LayerTypeSCTP,
	})
	// LayerClassIPControl contains TCP/IP control protocols.
	LayerClassIPControl = gopacket_dpdk.NewLayerClass([]gopacket_dpdk.LayerType{
		LayerTypeICMPv4,
		LayerTypeICMPv6,
	})
	// LayerClassSCTPChunk contains SCTP chunk types (not the top-level SCTP
	// layer).
	LayerClassSCTPChunk = gopacket_dpdk.NewLayerClass([]gopacket_dpdk.LayerType{
		LayerTypeSCTPUnknownChunkType,
		LayerTypeSCTPData,
		LayerTypeSCTPInit,
		LayerTypeSCTPSack,
		LayerTypeSCTPHeartbeat,
		LayerTypeSCTPError,
		LayerTypeSCTPShutdown,
		LayerTypeSCTPShutdownAck,
		LayerTypeSCTPCookieEcho,
		LayerTypeSCTPEmptyLayer,
		LayerTypeSCTPInitAck,
		LayerTypeSCTPHeartbeatAck,
		LayerTypeSCTPAbort,
		LayerTypeSCTPShutdownComplete,
		LayerTypeSCTPCookieAck,
	})
	// LayerClassIPv6Extension contains IPv6 extension headers.
	LayerClassIPv6Extension = gopacket_dpdk.NewLayerClass([]gopacket_dpdk.LayerType{
		LayerTypeIPv6HopByHop,
		LayerTypeIPv6Routing,
		LayerTypeIPv6Fragment,
		LayerTypeIPv6Destination,
	})
	LayerClassIPSec = gopacket_dpdk.NewLayerClass([]gopacket_dpdk.LayerType{
		LayerTypeIPSecAH,
		LayerTypeIPSecESP,
	})
)
