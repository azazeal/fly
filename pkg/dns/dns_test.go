package dns

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/azazeal/fly/internal/testutil"
)

func TestRegions(t *testing.T) {
	appName := token(t)

	defer stub(&mockResolver{
		lookupTXT: func(_ context.Context, name string) ([]string, error) {
			if name != "regions."+appName+".internal" {
				return nil, assert.AnError
			}

			return []string{
				"iad,cdg",
				"dfw,cdg,ewr",
				"ams,atl",
			}, nil
		},
	})()

	got, err := Regions(context.TODO(), appName)
	assert.NoError(t, err)
	assert.ElementsMatch(t, got, []string{
		"ams",
		"atl",
		"cdg",
		"cdg", // appears twice on purpose
		"dfw",
		"ewr",
		"iad",
	})
}

func TestAllInstances(t *testing.T) {
	appName := token(t)

	defer stub(&mockResolver{
		lookupIP: func(_ context.Context, network, name string) ([]net.IP, error) {
			if network != "ip6" || name != "global."+appName+".internal" {
				return nil, assert.AnError
			}

			return []net.IP{
				net.ParseIP("fdaa:0:22b7:a7b:ab8:3071:ecb3:2"),
				net.ParseIP("fdaa:0:22b7:a7b:abd:aa3c:6498:2"),
			}, nil
		},
	})()

	got, err := AllInstances(context.TODO(), appName)
	assert.NoError(t, err)
	assert.ElementsMatch(t, got, []net.IP{
		net.ParseIP("fdaa:0:22b7:a7b:abd:aa3c:6498:2"),
		net.ParseIP("fdaa:0:22b7:a7b:ab8:3071:ecb3:2"),
	})
}

func TestApps(t *testing.T) {
	defer stub(&mockResolver{
		lookupTXT: func(_ context.Context, name string) ([]string, error) {
			if name != "_apps.internal" {
				return nil, assert.AnError
			}

			return []string{
				"app2",
				"app3,app1",
				"app2,app4",
			}, nil
		},
	})()

	got, err := Apps(context.TODO())
	assert.NoError(t, err)
	assert.ElementsMatch(t, got, []string{
		"app1",
		"app2",
		"app2", // appears twice on purpose
		"app3",
		"app4",
	})
}

func TestPeers(t *testing.T) {
	defer stub(&mockResolver{
		lookupTXT: func(_ context.Context, name string) ([]string, error) {
			if name != "_peer.internal" {
				return nil, assert.AnError
			}

			return []string{
				"peer1,peer2,peer3",
				"peer4,peer1",
				"peer3,peer4",
			}, nil
		},
	})()

	got, err := Peers(context.TODO())
	assert.NoError(t, err)
	assert.ElementsMatch(t, got, []string{
		"peer1",
		"peer1", // appears twice on purpose
		"peer2",
		"peer3",
		"peer3", // appears twice on purpose
		"peer4",
		"peer4", // appears twice on purpose
	})
}

func TestPeer(t *testing.T) {
	const (
		ip1 = "fdaa:0:22b7:a8b:ce2:0:a:c02"
		ip2 = "fdaa:0:22b7:a7b:ce2:0:a:c02"
	)

	peer := token(t)

	defer stub(&mockResolver{
		lookupIP: func(_ context.Context, network, name string) ([]net.IP, error) {
			if network != "ip6" || name != peer+"._peer.internal" {
				return nil, assert.AnError
			}

			return []net.IP{
				testutil.ParseIP(t, ip2),
				testutil.ParseIP(t, ip1),
			}, nil
		},
	})()

	got, err := Peer(context.TODO(), peer)
	assert.NoError(t, err)
	assert.Equal(t, testutil.ParseIP(t, ip2), got)
}

type mockResolver struct {
	lookupTXT func(ctx context.Context, name string) ([]string, error)
	lookupIP  func(ctx context.Context, network, host string) ([]net.IP, error)
}

func (mr *mockResolver) LookupTXT(ctx context.Context, name string) ([]string, error) {
	return mr.lookupTXT(ctx, name)
}

func (mr *mockResolver) LookupIP(ctx context.Context, network, host string) ([]net.IP, error) {
	return mr.lookupIP(ctx, network, host)
}

func stub(r Resolver) func() {
	old := global
	global = New(r)

	return func() { global = old }
}

func token(t *testing.T) string {
	t.Helper()

	return testutil.HexString(t, 16)
}
