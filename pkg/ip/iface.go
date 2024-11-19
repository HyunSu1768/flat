package ip

import (
	"errors"
	"github.com/vishvananda/netlink"
	"net"
	"syscall"
)

func GetIfaceAddrs(iface *net.Interface) ([]netlink.Addr, error) {
	link := &netlink.Device{
		LinkAttrs: netlink.LinkAttrs{
			Index: iface.Index,
		},
	}
	return netlink.AddrList(link, syscall.AF_INET)
}

func GetIfaceV6Addrs(iface *net.Interface) ([]netlink.Addr, error) {
	link := &netlink.Device{
		LinkAttrs: netlink.LinkAttrs{
			Index: iface.Index,
		},
	}
	return netlink.AddrList(link, syscall.AF_INET6)
}

func GetInterfaceIP4Addrs(iface *net.Interface) ([]net.IP, error) {
	addrs, err := GetIfaceAddrs(iface)
	if err != nil {
		return nil, err
	}

	ipAddr := make([]net.IP, 0)

	ll := make([]net.IP, 0)

	for _, addr := range addrs {
		if addr.IP.To4() == nil {
			continue
		}

		if addr.IP.IsLinkLocalUnicast() {
			ll = append(ll, addr.IP)
		}

		if addr.IP.IsGlobalUnicast() {
			ipAddr = append(ipAddr, addr.IP)
		}
	}

	if len(ll) > 0 {
		ipAddr = append(ipAddr, ll...)
	}

	if len(ipAddr) > 0 {
		return ipAddr, nil
	}

	return nil, errors.New("Interface에 할당된 IP4 가 없습니다.")
}

func GetInterfaceIP6Addrs(iface *net.Interface) ([]net.IP, error) {
	addrs, err := GetIfaceV6Addrs(iface)
	if err != nil {
		return nil, err
	}

	ipAddr := make([]net.IP, 0)

	ll := make([]net.IP, 0)

	for _, addr := range addrs {
		if addr.IP.To4() == nil {
			continue
		}

		if addr.IP.IsLinkLocalUnicast() {
			ll = append(ll, addr.IP)
		}

		if addr.IP.IsGlobalUnicast() {
			ipAddr = append(ipAddr, addr.IP)
		}
	}

	if len(ll) > 0 {
		ipAddr = append(ipAddr, ll...)
	}

	if len(ipAddr) > 0 {
		return ipAddr, nil
	}

	return nil, errors.New("Interface에 할당된 IP6 가 없습니다.")
}

func GetInterfaceIP4AddrMatch(iface *net.Interface, matchIP net.IP) error {
	addrs, err := GetIfaceAddrs(iface)
	if err != nil {
		return err
	}

	for _, addr := range addrs {
		if addr.IP.To4() != nil {
			if addr.IP.To4().Equal(matchIP) {
				return nil
			}
		}
	}

	return errors.New("Interface에서 주어진 IP4와 일치하는 IP가 존재하지 않습니다.")
}

func GetInterfaceIP6AddrMatch(iface *net.Interface, matchIP net.IP) error {
	addrs, err := GetIfaceV6Addrs(iface)
	if err != nil {
		return err
	}

	for _, addr := range addrs {
		if addr.IP.To4() != nil {
			if addr.IP.To4().Equal(matchIP) {
				return nil
			}
		}
	}

	return errors.New("Interface에서 주어진 IP6와 일치하는 IP가 존재하지 않습니다.")
}
