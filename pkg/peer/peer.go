// Package peer implements functionality related to the current instance's peers.
package peer

import (
	"net"

	"github.com/azazeal/fly/pkg/env"
	"github.com/azazeal/fly/pkg/instance"
)

// IP returns the IP address of the named peer.
func IP(hostname string) (net.IP, error) {
	return instance.IP(hostname, env.AppName())
}
