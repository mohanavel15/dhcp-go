package main

import (
	"fmt"
	"net"
	"time"
)

type DHCPServer struct {
	config    *Config
	allocator *Allocator
}

func NewDHCPServer(config *Config, allocator *Allocator) DHCPServer {
	return DHCPServer{
		config:    config,
		allocator: allocator,
	}
}

func (ds *DHCPServer) Handle(m *Message) (Message, error) {
	response := Message{}

	if m.Op != BOOTREQUEST {
		return response, fmt.Errorf("not a boot request")
	}

	response.Op = BOOTREPLY
	response.Htype = m.Htype
	response.Hlen = m.Hlen
	response.Hops = m.Hops
	response.Xid = m.Xid
	response.Secs = m.Secs
	response.Flags = m.Flags

	response.Siaddr = ds.config.ServerIP
	response.Giaddr = net.IPv4zero
	response.Chaddr = m.Chaddr

	response.Sname = m.Sname
	response.File = m.File
	response.MagicCookie = m.MagicCookie

	response.Options = DhcpOpts{}
	response.Options.AddServerIP(ds.config.ServerIP)

	dhcp_message_type := DHCPMessageType(0)
	for _, option := range m.Options {
		if option.Type == byte(MessageType) {
			dhcp_message_type = DHCPMessageType(option.Data[0])
		}
	}

	var err error = nil

	switch dhcp_message_type {
	case DHCPDISCOVER:
		ds.handleDiscover(m, &response)
	case DHCPREQUEST:
		ds.handleRequest(m, &response)
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

	response.Options.AddEnd()

	return response, err
}

func (ds *DHCPServer) handleDiscover(m *Message, r *Message) {
	ip, err := ds.allocator.GetAvailableIP(m.Chaddr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("available ip???:", ip)

	r.Ciaddr = ip
	r.Yiaddr = ip

	options := DhcpOpts{}
	options.AddLeaseTime(5)
	options.AddMessageType(DHCPOFFER)

	r.Options = append(r.Options, options...)
}

func (ds *DHCPServer) handleRequest(m *Message, r *Message) {
	ip := net.IPv4zero

	for _, option := range m.Options {
		if option.Type == byte(RequestedIP) {
			ip = net.IP(option.Data)
			break
		}
	}

	network := net.IPNet{
		IP:   ds.config.NetworkID,
		Mask: ds.config.Subnet(),
	}

	if !network.Contains(ip) {
		r.Options.AddMessageType(DHCPNAK)
		return
	}

	err := ds.allocator.Allocate("", m.Chaddr, ip, time.Now().Add(time.Hour).Unix())
	if err != nil {
		r.Options.AddMessageType(DHCPNAK)
		return
	}

	r.Ciaddr = ip
	r.Yiaddr = ip

	r.Options.AddLeaseTime(60 * 60)
	r.Options.AddMessageType(DHCPACK)

	r.Options.AddNetmask(ds.config.Netmask)
	// r.Options.AddRouter(ds.config.RouterIP)
}
