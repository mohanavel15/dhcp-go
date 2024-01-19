# DHCP
Implementation of RFC 2131

[RFC 2131 - Dynamic Host Configuration Protocol](https://www.rfc-editor.org/rfc/rfc2131)

[RFC 1533 - DHCP Options and BOOTP Vendor Extensions](https://www.rfc-editor.org/rfc/rfc1533)

- [X] DHCP DISCOVER
- [X] DHCP OFFER
- [X] DHCP REQUEST
- [X] DHCP ACK
- [X] DHCP NAK
- [ ] DHCP DECLINE
- [ ] DHCP RELEASE
- [ ] DHCP INFORM

# How to use

dhcp-config.json
```json
{
	"interface": "eth0",
	"lease_time": 3600
}
```

Commands
```bash
$ ifconfig eth0 10.10.10.1/24 up
$ ip route add 10.10.10.0/24 dev eth0
$ go build
$ ./dhcp
```
