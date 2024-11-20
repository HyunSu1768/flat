package ip

import (
	"flat/pkg/ns"
	"github.com/vishvananda/netlink"
	"net"
	"syscall"
	"testing"
)

func TestEnsureV4AddressOnLink(t *testing.T) {
	teardown := ns.SetUpNetLinkTest(t)
	defer teardown()

	lo, err := netlink.LinkByName("lo")
	if err != nil {
		t.Fatal(err)
	}

	if netlink.LinkSetUp(lo); err != nil {
		t.Fatal(err)
	}

	ipn := IP4Net{IP: FromIP(net.ParseIP("127.0.0.2")), PrefixLen: 24}
	if err := EnsureV4AddressOnLink(ipn, ipn, lo); err != nil {
		t.Fatal(err)
	}

	addrs, err := netlink.AddrList(lo, syscall.AF_INET)
	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 1 || addrs[0].String() != "127.0.0.2/24 lo" {
		t.Fatalf("주소 %v는 예상되지 않습니다", addrs)
	}

	if err := netlink.AddrAdd(lo, &netlink.Addr{IPNet: &net.IPNet{IP: net.ParseIP("127.0.1.1"), Mask: net.CIDRMask(24, 32)}}); err != nil {
		t.Fatal(err)
	}

	if err := EnsureV4AddressOnLink(ipn, ipn, lo); err != nil {
		t.Fatal(err)
	}

	addrs, err = netlink.AddrList(lo, syscall.AF_INET)
	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 2 {
		t.Fatalf("2개의 주소가 예상되었습니다, addr : %v", addrs)
	}
}
