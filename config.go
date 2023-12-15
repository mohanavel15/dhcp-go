package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

type Config struct {
	Hostname  string `json:"hostname"`
	ServerIP  net.IP `json:"sever_ip"`
	RouterIP  net.IP `json:"router_ip"`
	NetworkID net.IP `json:"network_id"`
	Netmask   net.IP `json:"netmask"`
	Dns       net.IP `json:"dns"`
	LeaseTime uint32 `json:"lease_time"`
}

func (c *Config) Subnet() net.IPMask {
	return net.IPMask(c.Netmask.To4())
}

func LoadConfig(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var config Config

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func SaveConfig(path string, config Config) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(config)
	if err != nil {
		log.Fatal(err)
	}
}
