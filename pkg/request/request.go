// Package request implements requested-related functionality.
package request

import (
	"net"
	"net/http"
	"strconv"
)

// IDHeader denotes the header ID uses.
const IDHeader = "Fly-Request-Id"

// ID returns the ID of the request.
func ID(r *http.Request) string {
	return fetch(r, IDHeader)
}

// RegionHeader denotes the header Region uses.
const RegionHeader = "Fly-Region"

// Region returns the code of the region region at which the connection for r
// was accepted in and routed from.
func Region(r *http.Request) string {
	return fetch(r, RegionHeader)
}

// ClientIPHeader denotes the header ClientIP uses.
const ClientIPHeader = "Fly-Client-Ip"

// ClientIP returns the client IP address r carries.
func ClientIP(r *http.Request) net.IP {
	return net.ParseIP(fetch(r, ClientIPHeader))
}

// ForwardedPortHeader denotes the header ForwardedPort uses.
const ForwardedPortHeader = "Fly-Forwarded-Port"

// ForwardedPort returns the port of the edge node at which the request
// originated from
func ForwardedPort(r *http.Request) int {
	s := fetch(r, ForwardedPortHeader)
	if v, err := strconv.ParseUint(s, 10, 16); err == nil {
		return int(v)
	}

	return 0
}

func fetch(r *http.Request, key string) (val string) {
	if v, ok := r.Header[key]; ok && len(v) > 0 {
		val = v[0]
	}

	return
}
