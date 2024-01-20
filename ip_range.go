package main

import "net"

func GetIPs(ipNet *net.IPNet) ([]net.IP, error) {
	var ips []net.IP
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); inc(&ip) {
		ipCopy := make(net.IP, len(ip))
		copy(ipCopy, ip)
		ips = append(ips, ipCopy)
	}

	if len(ips) < 3 {
		return ips, nil
	}
	return ips[2 : len(ips)-1], nil
}

func inc(ip *net.IP) {
	for j := len(*ip) - 1; j >= 0; j-- {
		(*ip)[j]++
		if (*ip)[j] != 0 {
			break
		}
	}
}
