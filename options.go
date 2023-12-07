package main

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

type DHCPOptionType uint8

const (
	RequestedIPAddress   DHCPOptionType = 50
	IPAddressLeaseTime   DHCPOptionType = 51
	OptionOverload       DHCPOptionType = 52
	DHCPMessageTypeOP    DHCPOptionType = 53
	ServerIdentifier     DHCPOptionType = 54
	ParameterRequestList DHCPOptionType = 55
	MessageOP            DHCPOptionType = 56
	MaximumMessageSize   DHCPOptionType = 57
	ClientIdentifier     DHCPOptionType = 61
	End                  DHCPOptionType = 255
)

type DHCPOption struct {
	Type   uint8
	Length uint8
	Data   []uint8
}
