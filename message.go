package main

import (
	"encoding/hex"
	"net"
	"strconv"
)

type OP uint8

const (
	_ OP = iota
	BOOTREQUEST
	BOOTREPLY
)

type Message struct {
	Op          OP
	Htype       byte
	Hlen        byte
	Hops        byte
	Xid         string
	Secs        uint16
	Flags       uint16
	Ciaddr      net.IP
	Yiaddr      net.IP
	Siaddr      net.IP
	Giaddr      net.IP
	Chaddr      net.HardwareAddr
	Sname       [64]byte
	File        [128]byte
	MagicCookie [4]byte
	Options     DhcpOpts
}

func (m *Message) String() string {
	str := "op: "
	if m.Op == BOOTREQUEST {
		str += "boot request"
	} else if m.Op == BOOTREPLY {
		str += "boot reply"
	} else {
		str += "unknown"
	}

	str += "\nhtype: "
	str += strconv.FormatUint(uint64(m.Htype), 10)

	str += "\nhlen: "
	str += strconv.FormatUint(uint64(m.Hlen), 10)

	str += "\nhops: "
	str += strconv.FormatUint(uint64(m.Hops), 10)

	str += "\nxid: 0x"
	str += m.Xid

	str += "\nsecs: "
	str += strconv.FormatUint(uint64(m.Secs), 10)

	str += "\nflags: "
	str += strconv.FormatUint(uint64(m.Flags), 10)

	str += "\nclient ip address: "
	str += m.Ciaddr.String()

	str += "\nyour (client) ip address: "
	str += m.Yiaddr.String()

	str += "\nnext server ip address: "
	str += m.Siaddr.String()

	str += "\nrelay agent ip address: "
	str += m.Giaddr.String()

	str += "\nclient hardware address: "
	str += m.Chaddr.String()

	str += ""

	return str
}

func UnpackMessage(buf []byte) Message {
	m := Message{
		Op:          OP(buf[0]),
		Htype:       buf[1],
		Hlen:        buf[2],
		Hops:        buf[3],
		Xid:         hex.EncodeToString(buf[4:8]),
		Secs:        uint16(buf[8])<<8 | uint16(buf[9]),
		Flags:       uint16(buf[10])<<8 | uint16(buf[11]),
		Ciaddr:      buf[12:16],
		Yiaddr:      buf[16:20],
		Siaddr:      buf[20:24],
		Giaddr:      buf[24:28],
		Chaddr:      buf[28:34],
		Sname:       [64]byte(buf[44:108]),
		File:        [128]byte(buf[108:236]),
		MagicCookie: [4]byte(buf[236:240]),
		Options:     []DhcpOpt{},
	}

	idx := 240
	for idx < len(buf) {
		option := DhcpOpt{}
		option.Type = buf[idx]
		if option.Type == uint8(Opts_End) {
			break
		}

		option.Length = buf[idx+1]
		eop := idx + 2 + int(option.Length)

		option.Data = buf[idx+2 : eop]
		idx = eop
		m.Options = append(m.Options, option)
	}

	return m
}

func PackMessage(m Message) []byte {
	buf := []byte{}
	buf = append(buf, byte(m.Op))
	buf = append(buf, m.Htype)
	buf = append(buf, m.Hlen)
	buf = append(buf, m.Hops)

	xid, _ := hex.DecodeString(m.Xid)
	buf = append(buf, xid...)

	buf = append(buf, byte(m.Secs>>8))
	buf = append(buf, byte(m.Secs))

	buf = append(buf, byte(m.Flags>>8))
	buf = append(buf, byte(m.Flags))

	buf = append(buf, m.Ciaddr.To4()...)
	buf = append(buf, m.Yiaddr.To4()...)
	buf = append(buf, m.Siaddr.To4()...)
	buf = append(buf, m.Giaddr.To4()...)
	buf = append(buf, m.Chaddr...)
	buf = append(buf, make([]byte, 10)...)

	buf = append(buf, m.Sname[:]...)
	buf = append(buf, m.File[:]...)
	buf = append(buf, m.MagicCookie[:]...)

	for _, option := range m.Options {
		buf = append(buf, option.Type)
		if option.Length > 0 {
			buf = append(buf, option.Length)
			buf = append(buf, option.Data...)
		}
	}

	return buf
}
