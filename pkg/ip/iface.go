package ip

import (
	"errors"
	"fmt"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	log "k8s.io/klog/v2"
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

	return nil, errors.New("Interface에 할당된 IP4 가 없습니다")
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

	return nil, errors.New("Interface에 할당된 IP6 가 없습니다")
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

	return errors.New("Interface에서 주어진 IP4와 일치하는 IP가 존재하지 않습니다")
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

	return errors.New("Interface에서 주어진 IP6와 일치하는 IP가 존재하지 않습니다")
}

func GetDefaultGatewayInterface() (*net.Interface, error) {
	routes, err := netlink.RouteList(nil, syscall.AF_INET)
	if err != nil {
		return nil, err
	}

	for _, route := range routes {
		if route.Dst == nil || route.Dst.String() == "0.0.0.0/0" {
			if route.LinkIndex <= 0 {
				return nil, errors.New("Default Gateway를 찾았지만, Interface를 결정할 수 없습니다")
			}
			return net.InterfaceByIndex(route.LinkIndex)
		}
	}

	return nil, errors.New("Defluat Gateway를 찾을 수 없습니다")
}

func GetDefaultV6GatewayInterface() (*net.Interface, error) {
	routes, err := netlink.RouteList(nil, syscall.AF_INET6)
	if err != nil {
		return nil, err
	}

	for _, route := range routes {
		if route.Dst == nil || route.Dst.String() == "::/0" {
			if route.LinkIndex <= 0 {
				return nil, errors.New("Default V6 Gateway를 찾았지만, Interface를 결정할 수 없습니다")
			}
			return net.InterfaceByIndex(route.LinkIndex)
		}
	}

	return nil, errors.New("Defluat V6 Gateway를 찾을 수 없습니다")
}

func GetInterfaceByIP(ip net.IP) (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		err := GetInterfaceIP4AddrMatch(&iface, ip)
		if err == nil {
			return &iface, nil
		}
	}
	return nil, errors.New("주어진 IP에 맞는 인터페이스가 존재하지 않습니다")
}

func GetInterfaceByIP6(ip net.IP) (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		err := GetInterfaceIP6AddrMatch(&iface, ip)
		if err == nil {
			return &iface, nil
		}
	}
	return nil, errors.New("주어진 IP6에 맞는 인터페이스가 존재하지 않습니다")
}

func GetInterfaceBySpecificIPRouting(ip net.IP) (*net.Interface, net.IP, error) {
	routes, err := netlink.RouteGet(ip)
	if err != nil {
		return nil, nil, err
	}

	for _, route := range routes {
		iface, err := net.InterfaceByIndex(route.LinkIndex)
		if err != nil {
			return nil, nil, fmt.Errorf("interface를 찾을 수 없습니다 : %v", err)
		} else {
			return iface, route.Src, nil
		}
	}

	return nil, nil, errors.New("주어진 IP에 대한 interface를 찾을 수 없습니다")
}

func DirectRouting(ip net.IP) (bool, error) {
	routes, err := netlink.RouteGet(ip)
	if err != nil {
		return false, fmt.Errorf("%v에 대한 라우트를 확인할 수 없습니다 : %v", ip, err)
	}
	if len(routes) == 1 && routes[0].Gw == nil {
		return true, nil
	}
	return false, nil
}

func EnsureV4AddressOnLink(ipa IP4Net, ipn IP4Net, link netlink.Link) error {
	addr := netlink.Addr{IPNet: ipa.ToIPNet()}
	existingAddrs, err := netlink.AddrList(link, unix.AF_INET)
	if err != nil {
		return err
	}

	var hasAddr bool
	for _, existingAddr := range existingAddrs {
		if existingAddr.Equal(addr) {
			hasAddr = true
		}

		if ipn.Contains(FromIP(existingAddr.IP)) {
			if err := netlink.AddrDel(link, &existingAddr); err != nil {
				return fmt.Errorf("%s에서 IP주소 %s를 삭제하는데 실패하였습니다 : %s", link.Attrs().Name, existingAddr.String(), err)
			}
			log.Infof("%s에서 IP %s를 삭제하였습니다", link.Attrs().Name, existingAddr.String())
		}
	}

	if !hasAddr {
		if err := netlink.AddrAdd(link, &addr); err != nil {
			return fmt.Errorf("%s에 IP %s를 추가하는데 실패하였습니다 : %s", link.Attrs().Name, addr.String(), err)
		}
	}

	return nil
}

func EnsureV6AddressOnLink(ipa IP4Net, ipn IP4Net, link netlink.Link) error {
	addr := netlink.Addr{IPNet: ipa.ToIPNet()}
	existingAddrs, err := netlink.AddrList(link, unix.AF_INET6)
	if err != nil {
		return err
	}

	var hasAddr bool
	for _, existingAddr := range existingAddrs {
		if existingAddr.Equal(addr) {
			hasAddr = true
		}

		if ipn.Contains(FromIP(existingAddr.IP)) {
			if err := netlink.AddrDel(link, &existingAddr); err != nil {
				return fmt.Errorf("%s에서 IPv6 주소 %s를 삭제하는데 실패하였습니다 : %s", link.Attrs().Name, existingAddr.String(), err)
			}
			log.Infof("%s에서 IPv6 %s를 삭제하였습니다", link.Attrs().Name, existingAddr.String())
		}
	}

	if !hasAddr {
		if err := netlink.AddrAdd(link, &addr); err != nil {
			return fmt.Errorf("%s에 IPv6 %s를 추가하는데 실패하였습니다 : %s", link.Attrs().Name, addr.String(), err)
		}
	}

	return nil
}
