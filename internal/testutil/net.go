package testutil

import (
	"net"
	"testing"
)

// ParseIP parses s as an IP address, returning the result.
//
// If s is not a valid textual representation of an IP address, ParseIP will
// fail tb.
func ParseIP(tb testing.TB, s string) (ip net.IP) {
	tb.Helper()

	if ip = net.ParseIP(s); ip == nil {
		tb.Fatalf("failed parsing %q as a net.IP", s)
	}

	return
}
