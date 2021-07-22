// Package instance implements functionality applicable to the local fly.io
// instance.
package instance

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/azazeal/fly/pkg/env"
)

// IP returns the IP address corresponding to named instance of the given app.
func IP(hostname, app string) (ip net.IP, err error) {
	fqdn := fmt.Sprintf("%s.vm.%s.internal", hostname, app)

	var ips []net.IP
	if ips, err = net.LookupIP(fqdn); err == nil && len(ips) > 0 {
		ip = ips[0]
	}

	return
}

var (
	cachedIPMu sync.Mutex
	cachedIP   net.IP
)

// PrivateIP returns a copy of the local instance's IP address.
func PrivateIP() (ip net.IP, err error) {
	cachedIPMu.Lock()
	defer cachedIPMu.Unlock()

	if cachedIP != nil {
		ip = append(ip, cachedIP...)

		return
	}

	var hostname string
	if hostname, err = os.Hostname(); err != nil {
		return
	}

	if ip, err = IP(hostname, env.AppName()); err == nil {
		cachedIP = append(cachedIP, ip...)
	}

	return
}
