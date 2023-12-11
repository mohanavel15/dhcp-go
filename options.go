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
	NetMask              uint8 = 1
	RouterIP             uint8 = 3
	RequestedIP          uint8 = 50
	IPLeaseTime          uint8 = 51
	OptionOverload       uint8 = 52
	MessageType          uint8 = 53
	ServerIdentifier     uint8 = 54
	ParameterRequestList uint8 = 55
	MessageOP            uint8 = 56
	MaximumMessageSize   uint8 = 57
	ClientIdentifier     uint8 = 61
	End                  uint8 = 255
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
	*opts = append(*opts, NewDhcpOpt(IPLeaseTime, 4, byte(secs>>24), byte(secs>>16), byte(secs>>8), byte(secs)))
}

func (opts *DhcpOpts) AddMessageType(mt DHCPMessageType) {
	*opts = append(*opts, NewDhcpOpt(MessageType, 1, uint8(mt)))
}

func (opts *DhcpOpts) AddServerIP(ip net.IP) {
	*opts = append(*opts, NewDhcpOpt(ServerIdentifier, 4, ip.To4()...))
}

func (opts *DhcpOpts) AddEnd() {
	*opts = append(*opts, NewDhcpOpt(End, 0))
}

func (opts *DhcpOpts) AddNetmask(netmask net.IP) {
	*opts = append(*opts, NewDhcpOpt(NetMask, 4, netmask.To4()...))
}
func (opts *DhcpOpts) AddRouter(ip net.IP) {
	*opts = append(*opts, NewDhcpOpt(RouterIP, 4, ip.To4()...))
}
