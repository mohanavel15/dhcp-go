package main

import "net"

type DHCPMessageType uint8

const (
	_ DHCPMessageType = iota
	DHCPDISCOVER
	DHCPOFFER
	DHCPREQUEST
	DHCPDECLINE
	DHCPACK
	DHCPNAK
	DHCPRELEASE
	DHCPINFORM
)

const (
	Opts_NetMask              uint8 = 1
	Opts_RouterIP             uint8 = 3
	Opts_RequestedIP          uint8 = 50
	Opts_IPLeaseTime          uint8 = 51
	Opts_OptionOverload       uint8 = 52
	Opts_MessageType          uint8 = 53
	Opts_ServerIdentifier     uint8 = 54
	Opts_ParameterRequestList uint8 = 55
	Opts_MessageOP            uint8 = 56
	Opts_MaximumMessageSize   uint8 = 57
	Opts_ClientIdentifier     uint8 = 61
	Opts_End                  uint8 = 255
)

type DhcpOpt struct {
	Type   uint8
	Length uint8
	Data   []uint8
}

func NewDhcpOpt(t uint8, l uint8, d ...uint8) DhcpOpt {
	return DhcpOpt{
		Type:   t,
		Length: l,
		Data:   d,
	}
}

type DhcpOpts []DhcpOpt

func (opts *DhcpOpts) AddLeaseTime(secs uint32) {
	*opts = append(*opts, NewDhcpOpt(Opts_IPLeaseTime, 4, byte(secs>>24), byte(secs>>16), byte(secs>>8), byte(secs)))
}

func (opts *DhcpOpts) AddMessageType(mt DHCPMessageType) {
	*opts = append(*opts, NewDhcpOpt(Opts_MessageType, 1, uint8(mt)))
}

func (opts *DhcpOpts) AddServerIP(ip net.IP) {
	*opts = append(*opts, NewDhcpOpt(Opts_ServerIdentifier, 4, ip.To4()...))
}

func (opts *DhcpOpts) AddEnd() {
	*opts = append(*opts, NewDhcpOpt(Opts_End, 0))
}

func (opts *DhcpOpts) AddNetmask(netmask net.IPMask) {
	*opts = append(*opts, NewDhcpOpt(Opts_NetMask, 4, netmask...))
}
func (opts *DhcpOpts) AddRouter(ip net.IP) {
	*opts = append(*opts, NewDhcpOpt(Opts_RouterIP, 4, ip.To4()...))
}
