package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

type Config struct {
	Hostname  string     `json:"hostname"`
	Interface string     `json:"interface"`
	net       *net.IPNet `json:"-"`
	LeaseTime uint32     `json:"lease_time"`
}

func (c *Config) Init() {
	itf, err := net.InterfaceByName(c.Interface)
	if err != nil {
		log.Fatalf("Error Getting Interface: %s\n", err.Error())
	}

	addrs, err := itf.Addrs()
	if err != nil {
		log.Fatalln(err.Error())
	}

	c.net = addrs[0].(*net.IPNet)
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
