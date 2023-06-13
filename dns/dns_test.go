package dns

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"sort"
	"testing"

	"github.com/azazeal/fly/internal/testutil"
)

func TestRegions(t *testing.T) {
	appName := token(t)

	rng := testutil.RNG(t)

	t.Cleanup(stub(&mockResolver{
		lookupTXT: func(_ context.Context, name string) ([]string, error) {
			if exp := "regions." + appName + ".internal"; exp != name {
				return nil, fmt.Errorf("wrong name: want %q, have %q", exp, name)
			}

			ret := []string{
				"iad,cdg",
				"dfw,cdg,ewr",
				"ams,atl",
			}
			rng.Shuffle(len(ret), func(i, j int) { ret[i], ret[j] = ret[j], ret[i] })

			return ret, nil
		},
	}))

	got, err := Regions(context.TODO(), appName)
	testutil.AssertEqual(t, nil, err)

	sort.Strings(got)
	testutil.AssertEqual(t, []string{
		"ams",
		"atl",
		"cdg",
		"cdg", // appears twice on purpose
		"dfw",
		"ewr",
		"iad",
	}, got)
}

func TestAllInstances(t *testing.T) {
	appName := token(t)

	t.Cleanup(stub(&mockResolver{
		lookupIP: func(_ context.Context, network, name string) ([]net.IP, error) {
			if want := "ip6"; want != network {
				return nil, fmt.Errorf("wrong network: want %q, have %q", want, name)
			} else if want := "global." + appName + ".internal"; want != name {
				return nil, fmt.Errorf("wrong name: want %q, have %q", want, name)
			}

			ret := []net.IP{
				net.ParseIP("fdaa:0:22b7:a7b:abd:aa3c:6498:2"),
				net.ParseIP("fdaa:0:22b7:a7b:ab8:3071:ecb3:2"),
			}

			return ret, nil
		},
	}))

	got, err := Instances(context.TODO(), appName, "")
	testutil.AssertEqual(t, nil, err)

	sort.Slice(got, func(i, j int) bool {
		return bytes.Compare(got[i], got[j]) == -1
	})
	testutil.AssertEqual(t, []net.IP{
		net.ParseIP("fdaa:0:22b7:a7b:ab8:3071:ecb3:2"),
		net.ParseIP("fdaa:0:22b7:a7b:abd:aa3c:6498:2"),
	}, got)
}

func TestApps(t *testing.T) {
	t.Cleanup(stub(&mockResolver{
		lookupTXT: func(_ context.Context, name string) ([]string, error) {
			if want := "_apps.internal"; want != name {
				return nil, fmt.Errorf("wrong name: want %q, have %q", want, name)
			}

			return []string{
				"app2",
				"app3,app1",
				"app2,app4",
			}, nil
		},
	}))

	got, err := Apps(context.TODO())
	testutil.AssertEqual(t, nil, err)

	sort.Strings(got)
	testutil.AssertEqual(t, []string{
		"app1",
		"app2",
		"app2", // appears twice on purpose
		"app3",
		"app4",
	}, got)
}

func TestPeers(t *testing.T) {
	t.Cleanup(stub(&mockResolver{
		lookupTXT: func(_ context.Context, name string) ([]string, error) {
			if want := "_peer.internal"; want != name {
				return nil, fmt.Errorf("wrong name: want %q, have %q", want, name)
			}

			return []string{
				"peer1,peer2,peer3",
				"peer4,peer1",
				"peer3,peer4",
			}, nil
		},
	}))

	got, err := Peers(context.TODO())
	testutil.AssertEqual(t, nil, err)

	sort.Strings(got)
	testutil.AssertEqual(t, []string{
		"peer1",
		"peer1", // appears twice on purpose
		"peer2",
		"peer3",
		"peer3", // appears twice on purpose
		"peer4",
		"peer4", // appears twice on purpose
	}, got)
}

func TestPeer(t *testing.T) {
	const (
		ip1 = "fdaa:0:22b7:a8b:ce2:0:a:c02"
		ip2 = "fdaa:0:22b7:a7b:ce2:0:a:c02"
	)

	peer := token(t)

	t.Cleanup(stub(&mockResolver{
		lookupIP: func(_ context.Context, network, name string) ([]net.IP, error) {
			if want := "ip6"; want != network {
				return nil, fmt.Errorf("wrong network: want %q, have %q", want, name)
			} else if want := peer + "._peer.internal"; want != name {
				return nil, fmt.Errorf("wrong name: want %q, have %q", want, name)
			}

			return []net.IP{
				testutil.ParseIP(t, ip2),
				testutil.ParseIP(t, ip1),
			}, nil
		},
	}))

	got, err := Peer(context.TODO(), peer)
	testutil.AssertEqual(t, nil, err)
	testutil.AssertEqual(t, testutil.ParseIP(t, ip2), got)
}

func TestPrivateIP(t *testing.T) {
	const (
		ip1 = "fdaa:0:22b7:a7b:aa0:12a5:aacb:2"
		ip2 = "fdaa:0:22b7:a7b:ab8:244c:ae91:2"
	)

	t.Cleanup(stub(&mockResolver{
		lookupIP: func(_ context.Context, network, name string) ([]net.IP, error) {
			if want := "ip6"; want != network {
				return nil, fmt.Errorf("wrong network: want %q, have %q", want, name)
			} else if want := "fly-local-6pn"; want != name {
				return nil, fmt.Errorf("wrong name: want %q, have %q", want, name)
			}

			return []net.IP{
				testutil.ParseIP(t, ip2),
				testutil.ParseIP(t, ip1),
			}, nil
		},
	}))

	got, err := PrivateIP(context.TODO())
	testutil.AssertEqual(t, nil, err)
	testutil.AssertEqual(t, testutil.ParseIP(t, ip2), got)
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
