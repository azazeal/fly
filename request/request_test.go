package request

import (
	"net"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strconv"
	"testing"

	"github.com/azazeal/fly/internal/testutil"
)

func TestHeadersAreCanonical(t *testing.T) {
	headers := map[string]string{
		"IDHeader":            IDHeader,
		"RegionHeader":        RegionHeader,
		"ClientIPHeader":      ClientIPHeader,
		"ForwardedPortHeader": ForwardedPortHeader,
	}

	for name := range headers {
		got := headers[name]

		t.Run(name, func(t *testing.T) {
			exp := textproto.CanonicalMIMEHeaderKey(got)

			testutil.AssertEqual(t, exp, got)
		})
	}
}

func TestID(t *testing.T) {
	exp := testutil.HexString(t, 26)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("fly-request-id", exp) //nolint:canonicalheader // fly dox specify this header

	testutil.AssertEqual(t, exp, ID(req))
}

func TestRegion(t *testing.T) {
	exp := testutil.HexString(t, 3)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("fly-region", exp) //nolint:canonicalheader // fly dox specify this header

	testutil.AssertEqual(t, exp, Region(req))
}

func TestClientIP(t *testing.T) {
	const (
		ipv4 = "195.75.14.147"
		ipv6 = "6e4:dc69:40ac:ddd6:65b1:a798:d8ca:822"
	)
	cases := []struct {
		lit string
		exp net.IP
	}{
		0: {},
		1: {
			lit: "not an ip",
		},
		2: {
			lit: ipv4,
			exp: testutil.ParseIP(t, ipv4),
		},
		3: {
			lit: ipv6,
			exp: testutil.ParseIP(t, ipv6),
		},
	}

	for caseIndex := range cases {
		kase := cases[caseIndex]

		t.Run(strconv.Itoa(caseIndex), func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("fly-client-ip", kase.lit) //nolint:canonicalheader // fly dox specify this header

			testutil.AssertEqual(t, kase.exp, ClientIP(req))
		})
	}
}

func TestForwardedPort(t *testing.T) {
	cases := []struct {
		lit string
		exp int
	}{
		0: {},
		1: {
			lit: "not a port",
		},
		2: {
			lit: "0",
		},
		3: {
			lit: "1",
			exp: 1,
		},
		4: {
			lit: "65535",
			exp: 65535,
		},
		5: {
			lit: "65536",
			exp: 0,
		},
	}

	for caseIndex := range cases {
		kase := cases[caseIndex]

		t.Run(strconv.Itoa(caseIndex), func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("fly-forwarded-port", kase.lit) //nolint:canonicalheader // fly dox specify this header

			testutil.AssertEqual(t, kase.exp, ForwardedPort(req))
		})
	}
}
