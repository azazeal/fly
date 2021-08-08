// Package dns implements DNS-relatd functionality.
package dns

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
)

// Resolver wraps the functionality that instances of DNS rely on.
//
// All instances of net.Resolver implement Resolver.
type Resolver interface {
	// LookupTXT returns the DNS TXT records for the given domain name.
	LookupTXT(ctx context.Context, name string) ([]string, error)

	// LookupIP looks up host for the given network. It returns a slice of that
	// host's IP addresses of the type specified by network. network must be one
	// of "ip", "ip4" or "ip6".
	LookupIP(ctx context.Context, network, host string) ([]net.IP, error)
}

var global = New(net.DefaultResolver)

// New returns an instance of DNS that uses the given Resolver.
func New(r Resolver) DNS {
	return &wrapper{
		Resolver: r,
	}
}

// DNS wraps the functionality the .internal addresses of fly provide.
type DNS interface {
	// Regions returns the regions the named application is deployed to.
	Regions(ctx context.Context, appName string) ([]string, error)

	// Instances returns the IP addresses for the instances of the named
	// application in the given region.
	Instances(ctx context.Context, region, appName string) ([]net.IP, error)

	// Apps returns the applications running in the current organization.
	Apps(ctx context.Context) ([]string, error)

	// Peers returns the names of all wireguard peers.
	Peers(ctx context.Context) ([]string, error)

	// Peer returns the IPv6 address of the named wireguard peer.
	Peer(ctx context.Context, name string) (net.IP, error)

	// PrivateIP returns the IPv6 address of the local instance.
	PrivateIP(ctx context.Context) (net.IP, error)
}

type wrapper struct {
	Resolver

	privateIPMu sync.Mutex // protects privateIP
	privateIP   net.IP
}

func (w *wrapper) splitTXT(ctx context.Context, name string) (tokens []string, err error) {
	var txts []string
	if txts, err = w.LookupTXT(ctx, name); err == nil {
		for _, txt := range txts {
			tokens = append(tokens, strings.Split(txt, ",")...)
		}
	}

	return
}

func (w *wrapper) Regions(ctx context.Context, appName string) ([]string, error) {
	return w.splitTXT(ctx, "regions."+appName+".internal")
}

func (w *wrapper) Instances(ctx context.Context, region, appName string) ([]net.IP, error) {
	return w.LookupIP(ctx, "ip6", region+"."+appName+".internal")
}

func (w *wrapper) Apps(ctx context.Context) ([]string, error) {
	return w.splitTXT(ctx, "_apps.internal")
}

func (w *wrapper) Peers(ctx context.Context) ([]string, error) {
	return w.splitTXT(ctx, "_peer.internal")
}

func (w *wrapper) Peer(ctx context.Context, name string) (ip net.IP, err error) {
	host := fmt.Sprintf("%s._peer.internal", name)

	var ips []net.IP
	if ips, err = w.LookupIP(ctx, "ip6", host); err == nil && len(ips) > 0 {
		ip = ips[0]
	}

	return
}

func (w *wrapper) PrivateIP(ctx context.Context) (ip net.IP, err error) {
	w.privateIPMu.Lock()
	defer w.privateIPMu.Unlock()

	if w.privateIP == nil {
		const host = "fly-local-6pn"

		var ips []net.IP
		if ips, err = w.LookupIP(ctx, "ip6", host); err == nil && len(ips) > 0 {
			w.privateIP = append(w.privateIP, ips[0]...)
		}
	}

	ip = append(ip, w.privateIP...)

	return
}

// Regions returns the regions the named application is deployed to.
func Regions(ctx context.Context, appName string) ([]string, error) {
	return global.Regions(ctx, appName)
}

// AllInstances returns the IPv6 addresses for all of the instances of the
// named application.
func AllInstances(ctx context.Context, appName string) ([]net.IP, error) {
	return Instances(ctx, "global", appName)
}

// Instances returns the IP addresses for all of the instances of the named
// application in the given region.
func Instances(ctx context.Context, region, appName string) ([]net.IP, error) {
	return global.Instances(ctx, region, appName)
}

// Apps returns the applications running in the current organization.
func Apps(ctx context.Context) ([]string, error) {
	return global.Apps(ctx)
}

// Peers returns the names of all wireguard peers.
func Peers(ctx context.Context) ([]string, error) {
	return global.Peers(ctx)
}

// Peer returns the IPv6 address of the named wireguard peer.
func Peer(ctx context.Context, name string) (net.IP, error) {
	return global.Peer(ctx, name)
}

// PrivateIP returns the IPv6 address of the local instance.
func PrivateIP(ctx context.Context) (net.IP, error) {
	return global.PrivateIP(ctx)
}
