package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Client struct {
	Name    string
	MAC     net.HardwareAddr
	IP      net.IP
	Expires int64
}

type Allocator struct {
	AvailableIPs []net.IP
	Allocated    []Client
	mutex        sync.Mutex
}

func (a *Allocator) Clock() {
	ticker := time.NewTicker(time.Second)
	for {
		ct := <-ticker.C
		for i, c := range a.Allocated {
			if c.Expires < ct.Unix() {
				a.mutex.Lock()
				a.AvailableIPs = append(a.AvailableIPs, c.IP)
				a.Allocated = append(a.Allocated[:i], a.Allocated[i+1:]...)
				a.mutex.Unlock()
			}
		}
	}
}

func (a *Allocator) GetAvailableIP(mac net.HardwareAddr) (net.IP, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for _, c := range a.Allocated {
		if c.MAC.String() == mac.String() {
			return c.IP, nil
		}
	}

	if len(a.AvailableIPs) == 0 {
		return net.IPv4zero, fmt.Errorf("no more available IPs")
	}

	ip := a.AvailableIPs[0]
	a.AvailableIPs = a.AvailableIPs[1:]

	a.Allocated = append(a.Allocated, Client{
		MAC:     mac,
		IP:      ip,
		Expires: time.Now().Add(time.Second * 5).Unix(),
	})

	return ip, nil
}

func (a *Allocator) Release(mac net.HardwareAddr) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for i, c := range a.Allocated {
		if c.MAC.String() == mac.String() {
			a.Allocated = append(a.Allocated[:i], a.Allocated[i+1:]...)
			a.AvailableIPs = append(a.AvailableIPs, c.IP)
			return
		}
	}
}

func (a *Allocator) Renew(mac net.HardwareAddr, Expires int64) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for i, c := range a.Allocated {
		if c.MAC.String() == mac.String() {
			a.Allocated[i].Expires = Expires
			return
		}
	}
}

func (a *Allocator) RARP(mac net.HardwareAddr) (net.IP, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for _, c := range a.Allocated {
		if c.MAC.String() == mac.String() {
			return c.IP, nil
		}
	}

	return net.IPv4zero, fmt.Errorf("no IP found for MAC")
}

func (a *Allocator) Allocate(name string, mac net.HardwareAddr, ip net.IP, Expires int64) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if ip.IsUnspecified() {
		return fmt.Errorf("IP is unspecified")
	}

	for _, c := range a.Allocated {
		if c.IP.Equal(ip) && c.MAC.String() != mac.String() {
			return fmt.Errorf("IP already allocated")
		}
	}

	for i, c := range a.Allocated {
		if c.MAC.String() == mac.String() {
			a.Allocated[i].Name = name
			a.Allocated[i].IP = ip
			a.Allocated[i].Expires = Expires
			return nil
		}
	}

	client := Client{
		Name:    name,
		MAC:     mac,
		IP:      ip,
		Expires: Expires,
	}

	a.Allocated = append(a.Allocated, client)

	return nil
}

func NewAllocator(ips []net.IP) *Allocator {
	alloc := Allocator{
		AvailableIPs: ips,
		Allocated:    []Client{},
	}

	return &alloc
}
