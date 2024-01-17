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

	fmt.Println(config.net.IP, config.net.Mask)

	broadcast := make(net.IP, len(config.net.IP))
	copy(broadcast, config.net.IP) // For IPv4 prefix used by go

	for i := range config.net.Mask {
		broadcast[12+i] = config.net.IP[12+i] | ^config.net.Mask[i]
	}

	allocator := NewAllocator([]net.IP{
		net.IPv4(10, 10, 10, 2),
		net.IPv4(10, 10, 10, 3),
		net.IPv4(10, 10, 10, 4),
		net.IPv4(10, 10, 10, 5),
	})
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

		log.Println(res.Chaddr, "->", res.Ciaddr)
		resbuf := PackMessage(res)

		copy(addr.IP, broadcast)
		conn.WriteToUDP(resbuf, addr)
	}
}
