package main

import (
	"fmt"
	"net"
)

func handleMessage(m *Message) (Message, error) {
	if m.Op != BOOTREQUEST {
		return Message{}, fmt.Errorf("NOT BOOT REQUEST")
	}

	dhcp_message_type := DHCPMessageType(0)
	for _, option := range m.Options {
		if option.Type == byte(DHCPMessageTypeOP) {
			dhcp_message_type = DHCPMessageType(option.Data[0])
		}
	}

	var response Message
	var err error = nil

	switch dhcp_message_type {
	case DHCPDISCOVER:
		response = handleDiscover(m)
	case DHCPREQUEST:
		response = handleRequest(m)
	// case DHCPDECLINE:
	// 	fmt.Println("DHCPDECLINE")
	// case DHCPRELEASE:
	// 	fmt.Println("DHCPRELEASE")
	// case DHCPINFORM:
	// 	fmt.Println("DHCPINFORM")
	default:
		err = fmt.Errorf("unknown DHCP message type")
		fmt.Println("Unknown DHCP message type")
		fmt.Println(dhcp_message_type)
		fmt.Println(m.Options)
	}

	return response, err
}

func handleDiscover(m *Message) Message {
	response := Message{}

	response.Op = BOOTREPLY
	response.Htype = m.Htype
	response.Hlen = m.Hlen
	response.Hops = m.Hops
	response.Xid = m.Xid
	response.Secs = m.Secs
	response.Flags = m.Flags

	response.Ciaddr = net.IPv4(10, 10, 5, 2)
	response.Yiaddr = net.IPv4(10, 10, 5, 2)
	response.Siaddr = net.IPv4(10, 10, 5, 1)
	response.Giaddr = net.IPv4(0, 0, 0, 0)
	response.Chaddr = m.Chaddr

	response.Sname = m.Sname
	response.File = m.File
	response.MagicCookie = m.MagicCookie

	leaseTime := DHCPOption{
		Type:   uint8(IPAddressLeaseTime),
		Length: 4,
		Data:   []uint8{},
	}

	var hour uint32 = 1 * 5 * 60
	leaseTime.Data = append(leaseTime.Data, byte(hour>>24))
	leaseTime.Data = append(leaseTime.Data, byte(hour>>16))
	leaseTime.Data = append(leaseTime.Data, byte(hour>>8))
	leaseTime.Data = append(leaseTime.Data, byte(hour))

	messageType := DHCPOption{
		Type:   uint8(DHCPMessageTypeOP),
		Length: 1,
		Data:   []uint8{uint8(DHCPOFFER)},
	}

	serverIP := DHCPOption{
		Type:   uint8(ServerIdentifier),
		Length: 4,
		Data:   net.IPv4(10, 10, 5, 1)[12:],
	}

	end := DHCPOption{
		Type:   uint8(End),
		Length: 0,
		Data:   []uint8{},
	}

	response.Options = append(response.Options, leaseTime, messageType, serverIP, end)

	return response
}

func handleRequest(m *Message) Message {
	response := Message{}

	response.Op = BOOTREPLY
	response.Htype = m.Htype
	response.Hlen = m.Hlen
	response.Hops = m.Hops
	response.Xid = m.Xid
	response.Secs = m.Secs
	response.Flags = m.Flags

	response.Ciaddr = net.IPv4(10, 10, 5, 2)
	response.Yiaddr = net.IPv4(10, 10, 5, 2)
	response.Siaddr = net.IPv4(10, 10, 5, 1)
	response.Giaddr = net.IPv4(0, 0, 0, 0)
	response.Chaddr = m.Chaddr

	response.Sname = m.Sname
	response.File = m.File
	response.MagicCookie = m.MagicCookie

	leaseTime := DHCPOption{
		Type:   uint8(IPAddressLeaseTime),
		Length: 4,
		Data:   []uint8{},
	}

	var hour uint32 = 1 * 5 * 60
	leaseTime.Data = append(leaseTime.Data, byte(hour>>24))
	leaseTime.Data = append(leaseTime.Data, byte(hour>>16))
	leaseTime.Data = append(leaseTime.Data, byte(hour>>8))
	leaseTime.Data = append(leaseTime.Data, byte(hour))

	messageType := DHCPOption{
		Type:   uint8(DHCPMessageTypeOP),
		Length: 1,
		Data:   []uint8{uint8(DHCPACK)},
	}

	serverIP := DHCPOption{
		Type:   uint8(ServerIdentifier),
		Length: 4,
		Data:   net.IPv4(10, 10, 5, 1)[12:],
	}

	end := DHCPOption{
		Type:   uint8(End),
		Length: 0,
		Data:   []uint8{},
	}

	response.Options = append(response.Options, leaseTime, messageType, serverIP, end)

	return response
}
