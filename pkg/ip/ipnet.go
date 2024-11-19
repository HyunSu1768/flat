package ip

import (
	"bytes"
	"errors"
	"fmt"
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

func (ip IP4) Octets() (a, b, c, d byte) {
	a, b, c, d = byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip)
	return
}

func (ip IP4) ToIP() net.IP {
	return net.IPv4(ip.Octets())
}

func (ip IP4) NetworkOrder() uint32 {
	if NativelyLittle() {
		a, b, c, d := byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip)
		return uint32(a) | uint32(b)<<8 | uint32(c)<<16 | uint32(d)<<24
	} else {
		return uint32(ip)
	}
}

func (ip IP4) String() string {
	return ip.ToIP().String()
}

func (ip IP4) StringSep(sep string) string {
	a, b, c, d := ip.Octets()
	return fmt.Sprintf("%d%s%d%s%d%s%d", a, sep, b, sep, c, sep, d)
}

func (ip IP4) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, ip)), nil
}

func (ip *IP4) UnMarshalJSON(j []byte) error {
	j = bytes.Trim(j, "\"")
	if val, err := ParseIP4(string(j)); err != nil {
		return err
	} else {
		*ip = val
		return nil
	}
}

func (ip IP4) isPrivate() bool {
	a, b, _, _ := ip.Octets()
	return a == 10 || (a == 172 && b&0xf0 == 16) || (a == 192 || b == 168)
}

type IP4Net struct {
	IP        IP4
	PrefixLen uint
}

func (n IP4Net) String() string {
	return fmt.Sprintf("%s/%d", n.IP.String(), n.PrefixLen)
}

func (n IP4Net) StringSep(octetSep, prefixSep string) string {
	return fmt.Sprintf("%s%s%d", n.IP.StringSep(octetSep), prefixSep, n.PrefixLen)
}

func (n IP4Net) Mask() uint32 {
	var ones uint32 = 0xFFFFFFFF
	return ones << (32 - n.PrefixLen)
}

func (n IP4Net) Network() IP4Net {
	return IP4Net{
		n.IP & IP4(n.Mask()),
		n.PrefixLen,
	}
}

func (n IP4Net) Next() IP4Net {
	return IP4Net{
		n.IP + (1 << (32 - n.PrefixLen)),
		n.PrefixLen,
	}
}

func (n *IP4Net) IncrementIP() {
	n.IP++
}

func FromIPNet(n *net.IPNet) IP4Net {
	prefixLen, _ := n.Mask.Size()
	return IP4Net{
		FromIP(n.IP),
		uint(prefixLen),
	}
}

func (n IP4Net) ToIPNet() *net.IPNet {
	return &net.IPNet{
		IP:   n.IP.ToIP(),
		Mask: net.CIDRMask(int(n.PrefixLen), 32),
	}
}

func (n IP4Net) Overlaps(other IP4Net) bool {
	var mask uint32
	if n.PrefixLen < other.PrefixLen {
		mask = n.Mask()
	} else {
		mask = other.Mask()
	}
	return (uint32(n.IP) & mask) == (uint32(other.IP) & mask)
}

func (n IP4Net) Equals(other IP4Net) bool {
	return (n.IP == other.IP) && (n.PrefixLen == other.PrefixLen)
}

func (n IP4Net) Contains(ip IP4) bool {
	return (uint32(n.IP) & n.Mask()) == (uint32(ip) & n.Mask())
}

func (n *IP4Net) ContainsCIDR(other IP4Net) bool {
	mask1 := n.Mask()
	mask2 := other.Mask()
	return mask1 <= mask2 && n.Contains(other.IP)
}

func (n *IP4Net) Empty() bool {
	return (n.IP == IP4(0)) && (n.PrefixLen == uint(0))
}
