package ip

import (
	"encoding/json"
	"net"
	"testing"
)

func makeIP4Net(s string, cidr uint) IP4Net {
	ip, err := ParseIP4(s)
	if err != nil {
		panic(err)
	}
	return IP4Net{
		ip,
		cidr,
	}
}

func mkIP4(s string) IP4 {
	ip, err := ParseIP4(s)
	if err != nil {
		panic(err)
	}
	return ip
}

func TestIP4(t *testing.T) {
	nip := net.ParseIP("1.2.3.4")
	ip := FromIP(nip)
	a, b, c, d := ip.Octets()
	if a != 1 || b != 2 || c != 3 || d != 4 {
		t.Error("FromIP 실패")
	}

	ip, err := ParseIP4("1.2.3.4")
	if err != nil {
		t.Error("PaseIP4 실패 : ", err)
	} else {
		a, b, c, d := ip.Octets()
		if a != 1 || b != 2 || c != 3 || d != 4 {
			t.Error("ParseIP4 실패")
		}
	}

	if ip.ToIP().String() != "1.2.3.4" {
		t.Error("ToIP 실패")
	}

	if ip.String() != "1.2.3.4" {
		t.Error("String 실패")
	}

	if ip.StringSep("-") != "1-2-3-4" {
		t.Error("StringSep 실패")
	}

	j, err := json.Marshal(ip)
	if err != nil {
		t.Error("IP4 Marshal 실패 : ", err)
	} else if string(j) != `"1.2.3.4"` {
		t.Error("IP4 Marshal 이 예기치 못한 이유로 실패 : ", err)
	}

	addresses := []*struct {
		ip      string
		private bool
	}{
		{"192.168.0.1", true},
		{"172.16.0.1", true},
		{"172.31.0.1", true},
		{"10.1.2.3", true},

		{"8.8.8.8", false},
		{"172.32.0.1", false},
		{"192.167.0.1", false},
		{"192.169.0.1", false},
	}

	for _, address := range addresses {
		ip := mkIP4(address.ip)
		is_private := ip.isPrivate()
		if is_private != address.private {
			t.Errorf("%v - 예상된 private: %v, 실제 private: %v", address.ip, address.private, is_private)
		}
	}
}

func TestIP4Net(t *testing.T) {
	n1 := makeIP4Net("1.2.3.0", 24)

	if n1.ToIPNet().String() != "1.2.3.0/24" {
		t.Error("ToIPNet 실패")
	}

	if n1.StringSep("*", "/") != "1*2*3*0/24" {
		t.Error("StringSep 실패")
	}

	if n1.Mask() != 0xFFFFFF00 {
		t.Error("Mask 실패")
	}

	if !n1.Overlaps(n1) {
		t.Errorf("%s은 %s를 Overlaps 하지 않는다.", n1, n1)
	}

	n2 := makeIP4Net("1.2.0.0", 16)
	if !n1.Overlaps(n2) {
		t.Errorf("%s은 %s를 Overlaps 하지 않는다.", n1, n2)
	}

	n2 = makeIP4Net("7.2.4.0", 22)
	if n1.Overlaps(n2) {
		t.Errorf("%s은 %s을 Overlaps 한다.", n1, n2)
	}

	if !n1.Contains(mkIP4("1.2.3.4")) {
		t.Error("Contains 실패")
	}

	if !n1.Contains(mkIP4("1.2.3.0")) {
		t.Error("Contains 실패")
	}

	if n1.Contains(mkIP4("2.3.4.5")) {
		t.Error("Contains 실패")
	}

	j, err := json.Marshal(n1)
	if err != nil {
		t.Error("IP4NET Marshal 실패 : ", err)
	} else if string(j) != `"1.2.3.0/24"` {
		t.Error("예기치 못한 이유로 IP4NET Marshal 실패 : ", err)
	}

	n1.IncrementIP()
	if n1.String() != "1.2.3.1/24" {
		t.Error("IncrementIP 실패")
	}
}
