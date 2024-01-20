package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	config := LoadConfig("./dhcp-config.json")
	config.Init()

	broadcast := make(net.IP, len(config.net.IP))
	copy(broadcast, config.net.IP) // For IPv4 prefix used by go

	for i := range config.net.Mask {
		broadcast[12+i] = config.net.IP[12+i] | ^config.net.Mask[i]
	}

	ips, err := GetIPs(config.net)
	if err != nil {
		fmt.Println("Error getting IPs:", err)
		os.Exit(1)
	}

	fmt.Println(ips)

	allocator := NewAllocator(ips)
	go allocator.Clock()

	DhcpServer := NewDHCPServer(&config, allocator)

	addr := net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 67,
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error listening on UDP:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Server is listening on port:", addr.String())

	for {
		buffer := make([]byte, 576)
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading UDP packet:", err)
			continue
		}

		msg := UnpackMessage(buffer[:n])

		res, err := DhcpServer.Handle(&msg)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		logmsg := fmt.Sprintf("IP: %s is allocated to %s", res.Ciaddr.String(), res.Chaddr.String())
		if msg.Options.GetHostName() != "" {
			logmsg += fmt.Sprintf("(%s)", msg.Options.GetHostName())
		}
		logmsg += fmt.Sprintf(" for %d seconds", res.Options.GetLeaseTime())
		log.Println(logmsg)

		resbuf := PackMessage(res)

		copy(addr.IP, broadcast)
		conn.WriteToUDP(resbuf, addr)
	}
}
