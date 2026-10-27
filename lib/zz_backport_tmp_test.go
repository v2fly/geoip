package lib

import (
	"net"
	"net/netip"
	"testing"
)

func TestTmpCommentLine(t *testing.T) {
	e := NewEntry("x")
	if err := e.AddPrefix("# comment"); err != nil {
		t.Errorf("comment line AddPrefix err: %v", err)
	}
	if err := e.RemovePrefix("// comment"); err != nil {
		t.Errorf("comment line RemovePrefix err: %v", err)
	}
}

func TestTmpStaleCache(t *testing.T) {
	e := NewEntry("x")
	_ = e.AddPrefix("1.0.0.0/24")
	got, _ := e.MarshalText()
	_ = e.AddPrefix("2.0.0.0/24")
	got, _ = e.MarshalText()
	if len(got) != 2 {
		t.Errorf("stale cache after add: %v", got)
	}
	_ = e.RemovePrefix("1.0.0.0/24")
	got, _ = e.MarshalText()
	if len(got) != 1 {
		t.Errorf("stale cache after remove: %v", got)
	}
}

func TestTmpContainerStaleCache(t *testing.T) {
	c := NewContainer()
	a := NewEntry("x")
	_ = a.AddPrefix("1.0.0.0/24")
	_ = c.Add(a)
	v, _ := c.GetEntry("x")
	_, _ = v.MarshalText()
	b := NewEntry("x")
	_ = b.AddPrefix("2.0.0.0/24")
	_ = b.AddPrefix("2001:db8::/32")
	_ = c.Add(b)
	got, _ := v.MarshalText()
	if len(got) != 3 {
		t.Errorf("container stale cache after add: %v", got)
	}
	r := NewEntry("x")
	_ = r.AddPrefix("1.0.0.0/24")
	_ = c.Remove(r, CaseRemovePrefix)
	got, _ = v.MarshalText()
	if len(got) != 2 {
		t.Errorf("container stale cache after remove: %v", got)
	}
}

func TestTmpMappedCIDR(t *testing.T) {
	cases := map[string]string{
		"::ffff:1.2.3.4/120": "1.2.3.0/24",
		"1.2.3.4/24":         "1.2.3.0/24",
		"2001:db8::1/32":     "2001:db8::/32",
		"1.2.3.4":            "1.2.3.4/32",
		"::ffff:1.2.3.4":     "1.2.3.4/32",
	}
	for in, want := range cases {
		e := NewEntry("x")
		if err := e.AddPrefix(in); err != nil {
			t.Errorf("%s: err %v", in, err)
			continue
		}
		got, err := e.MarshalText()
		if err != nil || len(got) != 1 || got[0] != want {
			t.Errorf("%s: got %v err %v, want %s", in, got, err, want)
		}
	}
	for _, bad := range []string{"::ffff:1.2.3.4/64", "1.2.3.4/33", "abc/24", "1.2.3.4/"} {
		e := NewEntry("x")
		if err := e.AddPrefix(bad); err == nil {
			got, _ := e.MarshalText()
			t.Errorf("%s: expected error, got %v", bad, got)
		}
	}
}

func TestTmpIPNet(t *testing.T) {
	cases := []struct {
		n    *net.IPNet
		want string
	}{
		{&net.IPNet{IP: net.ParseIP("1.2.3.0"), Mask: net.CIDRMask(24, 32)}, "1.2.3.0/24"},
		{&net.IPNet{IP: net.ParseIP("1.2.3.0"), Mask: net.CIDRMask(120, 128)}, "1.2.3.0/24"},
		{&net.IPNet{IP: net.ParseIP("1.2.3.0").To4(), Mask: net.CIDRMask(24, 32)}, "1.2.3.0/24"},
		{&net.IPNet{IP: net.ParseIP("2001:db8::"), Mask: net.CIDRMask(32, 128)}, "2001:db8::/32"},
	}
	for _, c := range cases {
		e := NewEntry("x")
		if err := e.AddPrefix(c.n); err != nil {
			t.Errorf("%v: err %v", c.n, err)
			continue
		}
		got, err := e.MarshalText()
		if err != nil || len(got) != 1 || got[0] != c.want {
			t.Errorf("%v: got %v err %v want %s", c.n, got, err, c.want)
		}
	}
}

func TestTmpNilPointers(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic: %v", r)
		}
	}()
	e := NewEntry("x")
	var a *netip.Addr
	var p *netip.Prefix
	var n *net.IPNet
	if err := e.AddPrefix(a); err == nil {
		t.Error("nil *netip.Addr: expected error")
	}
	if err := e.AddPrefix(p); err == nil {
		t.Error("nil *netip.Prefix: expected error")
	}
	if err := e.AddPrefix(n); err == nil {
		t.Error("nil *net.IPNet: expected error")
	}
}

func TestTmpAddrNotMutated(t *testing.T) {
	a := netip.MustParseAddr("::ffff:1.2.3.4")
	e := NewEntry("x")
	_ = e.AddPrefix(&a)
	if !a.Is4In6() {
		t.Errorf("caller's addr mutated: %v", a)
	}
}

func TestTmpCutterLoop(t *testing.T) {
	c := NewContainer()
	for _, n := range []string{"a", "b", "c", "d", "e", "f", "g", "h"} {
		e := NewEntry(n)
		_ = e.AddPrefix("1.0.0.0/24")
		_ = c.Add(e)
	}
	for e := range c.Loop() {
		_ = c.Remove(e, CaseRemoveEntry)
	}
	if c.Len() != 0 {
		t.Errorf("len %d", c.Len())
	}
}
