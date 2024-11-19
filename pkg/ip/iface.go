package ip

import (
	"github.com/vishvananda/netlink"
	"net"
	"syscall"
)

func GetIfaceAddr(iface *net.Interface) ([]netlink.Addr, error) {
	link := &netlink.Device{
		LinkAttrs: netlink.LinkAttrs{
			Index: iface.Index,
		},
	}
	return netlink.AddrList(link, syscall.AF_INET)
}
