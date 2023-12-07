package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
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

	buffer := make([]byte, 576)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading UDP packet:", err)
			continue
		}

		fmt.Printf("Received %d bytes from %s\n", n, addr.String())
		msg := UnpackMessage(buffer[:n])
		fmt.Println(msg.String())

		res, err := handleMessage(&msg)
		if err != nil {
			continue
		}
		resbuf := PackMessage(res)

		fmt.Println("Sending response...")

		ToIP := addr.IP
		if ToIP.Equal(net.IPv4zero) {
			ToIP = net.IPv4bcast
		}

		conn.WriteToUDP(resbuf, &net.UDPAddr{IP: ToIP, Port: 68})
	}
}
