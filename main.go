package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	config := LoadConfig("./dhcp-config.json")

	allocator := NewAllocator([]net.IP{
		net.IPv4(10, 10, 10, 2),
		net.IPv4(10, 10, 10, 3),
		net.IPv4(10, 10, 10, 4),
		net.IPv4(10, 10, 10, 5),
	})
	go allocator.Clock()

	DhcpServer := NewDHCPServer(&config, allocator)

	addr, err := net.ResolveUDPAddr("udp", ":67")
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening on UDP:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("UDP server is listening on port 67...")

	for {
		buffer := make([]byte, 576)
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading UDP packet:", err)
			continue
		}

		msg := UnpackMessage(buffer[:n])
		// fmt.Println(msg.String())

		res, err := DhcpServer.Handle(&msg)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		log.Printf("%s allocated to %s for %ds\n", res.Yiaddr.String(), res.Chaddr.String(), config.LeaseTime)
		resbuf := PackMessage(res)

		if addr.IP.IsUnspecified() {
			addr.IP = net.IPv4(10, 255, 255, 255)
		}

		conn.WriteToUDP(resbuf, addr)
	}
}
