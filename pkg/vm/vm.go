// Package vm implements virtual machine related functionality.
package vm

import (
	"fmt"
	"net"
	"sync"
)

// IP returns the IP address corresponding to the virtual machine of the given
// hostname for the given app.
func IP(hostname, app string) (net.IP, error) {
	fqdn := fmt.Sprintf("%s.vm.%s.internal", hostname, app)

	return lookupIP(fqdn)
}

var (
	cachedIPMu sync.Mutex
	cachedIP   net.IP
)

// PrivateIP returns a copy of the local virtual machine's IP address.
func PrivateIP() (ip net.IP, err error) {
	cachedIPMu.Lock()
	defer cachedIPMu.Unlock()

	if cachedIP != nil {
		ip = append(ip, cachedIP...)

		return
	}

	if ip, err = lookupIP("fly-local-6pn"); err == nil {
		cachedIP = append(cachedIP, ip...)
	}

	return
}

func lookupIP(host string) (ip net.IP, err error) {
	var ips []net.IP
	if ips, err = net.LookupIP(host); err == nil && len(ips) > 0 {
		ip = ips[0]
	}

	return
}
