package request

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/azazeal/fly/internal/testutil"
)

func TestHeadersAreCanonical(t *testing.T) {
	testutil.AssertHeadersAreCanonical(t, map[string]string{
		"IDHeader":            IDHeader,
		"RegionHeader":        RegionHeader,
		"ClientIPHeader":      ClientIPHeader,
		"ForwardedPortHeader": ForwardedPortHeader,
	})
}

func TestID(t *testing.T) {
	exp := testutil.HexString(t, 26)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("fly-request-id", exp)

	assert.Equal(t, exp, ID(req))
}

func TestRegion(t *testing.T) {
	exp := testutil.HexString(t, 3)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("fly-region", exp)

	assert.Equal(t, exp, Region(req))
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
			exp: parseIP(t, ipv4),
		},
		3: {
			lit: ipv6,
			exp: parseIP(t, ipv6),
		},
	}

	for caseIndex := range cases {
		kase := cases[caseIndex]

		t.Run(strconv.Itoa(caseIndex), func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("fly-client-ip", kase.lit)

			assert.Equal(t, kase.exp, ClientIP(req))
		})
	}
}

func parseIP(t *testing.T, s string) (ip net.IP) {
	t.Helper()

	ip = net.ParseIP(s)
	require.NotNil(t, ip)

	return
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
			req.Header.Set("fly-forwarded-port", kase.lit)

			assert.Equal(t, kase.exp, ForwardedPort(req))
		})
	}
}
