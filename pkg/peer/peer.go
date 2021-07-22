// Package peer implements peer-related functionality.
package peer

import (
	"net"

	"github.com/azazeal/fly/pkg/env"
	"github.com/azazeal/fly/pkg/vm"
)

// IPs returns the IP addresses of the virtual machine's peers.
func IPs() ([]net.IP, error) {
	return peers("global", env.AppName())
}

// RegionalIPs returns the IP addresses of the current virtual machine's
// regional peers.
func RegionalIPs() ([]net.IP, error) {
	return peers(env.Region(), env.AppName())
}

func peers(region, app string) (ips []net.IP, err error) {
	if ips, err = vm.IPs(region, app); err != nil {
		return
	}

	// remove our own ip from the slice
	var ip net.IP
	if ip, err = vm.PrivateIP(); err == nil {
		ips = excludeIP(ips, ip)
	}

	return
}

func excludeIP(slice []net.IP, ip net.IP) []net.IP {
	i := 0
	for _, n := range slice {
		if ip.Equal(n) {
			continue
		}

		slice[i] = n
		i++
	}

	return slice[:i]
}
