package ip

import (
	"errors"
	"net"
)

type IP4 uint32

func FromBytes(ip []byte) IP4 {
	return IP4(uint32(ip[3]) |
		(uint32(ip[2]) << 8) |
		(uint32(ip[1]) << 16) |
		(uint32(ip[0]) << 24))
}

func FromIP(ip net.IP) IP4 {
	ipv4 := ip.To4()

	if ipv4 == nil {
		panic("주소가 ipv4가 아닙니다")
	}

	return FromBytes(ip)
}

func ParseIP4(s string) (IP4, error) {
	ip := net.ParseIP(s)

	if ip == nil {
		return IP4(0), errors.New("올바르지 않은 IP 포멧입니다")
	}
	return FromIP(ip), nil
}

func MustParseIP4(s string) IP4 {
	ipv4, err := ParseIP4(s)
	if err != nil {
		panic(err)
	}
	return ipv4
}
