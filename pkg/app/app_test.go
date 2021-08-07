package app

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/azazeal/fly/internal/testutil"
	"github.com/azazeal/fly/pkg/env"
)

func TestList(t *testing.T) {
	apps := []string{
		testutil.HexString(t, 10),
		testutil.HexString(t, 10),
	}

	cases := []struct {
		lookupTXT func(context.Context, string) ([]string, error)
		exp       []string
		err       error
	}{
		0: {
			lookupTXT: func(context.Context, string) ([]string, error) {
				return nil, nil
			},
		},
		1: {
			lookupTXT: func(context.Context, string) ([]string, error) {
				return nil, assert.AnError
			},
			err: assert.AnError,
		},
		2: {
			lookupTXT: func(context.Context, string) ([]string, error) {
				return nil, nil
			},
		},
		3: {
			lookupTXT: func(_ context.Context, name string) ([]string, error) {
				if name != "_apps.internal" {
					panic(fmt.Errorf("wrong name: %q", name))
				}

				return apps, nil
			},
			exp: apps,
		},
	}

	for caseIndex := range cases {
		kase := cases[caseIndex]

		t.Run(strconv.Itoa(caseIndex), func(t *testing.T) {
			if kase.lookupTXT != nil {
				defer withLookupTXT(kase.lookupTXT)()
			}

			got, err := List(context.TODO())
			if kase.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Same(t, kase.err, err)
			}
			assert.EqualValues(t, kase.exp, got)
		})
	}
}

func TestRegions(t *testing.T) {
	appName := testutil.HexString(t, 5)
	defer testutil.SetEnv(t, map[string]string{
		env.AppNameKey: appName,
	})()

	regions := []string{
		testutil.HexString(t, 10),
		testutil.HexString(t, 10),
	}

	cases := []struct {
		lookupTXT func(context.Context, string) ([]string, error)
		exp       []string
		err       error
	}{
		0: {
			lookupTXT: func(context.Context, string) ([]string, error) {
				return nil, nil
			},
		},
		1: {
			lookupTXT: func(context.Context, string) ([]string, error) {
				return nil, assert.AnError
			},
			err: assert.AnError,
		},
		2: {
			lookupTXT: func(context.Context, string) ([]string, error) {
				return nil, nil
			},
		},
		3: {
			lookupTXT: func(_ context.Context, name string) ([]string, error) {
				if name != "regions."+appName+".internal" {
					panic(fmt.Errorf("wrong name: %q", name))
				}

				return regions, nil
			},
			exp: regions,
		},
	}

	for caseIndex := range cases {
		kase := cases[caseIndex]

		t.Run(strconv.Itoa(caseIndex), func(t *testing.T) {
			if kase.lookupTXT != nil {
				defer withLookupTXT(kase.lookupTXT)()
			}

			got, err := Regions(context.TODO())
			if kase.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Same(t, kase.err, err)
			}
			assert.EqualValues(t, kase.exp, got)
		})
	}
}

func withLookupTXT(fn func(context.Context, string) ([]string, error)) func() {
	old := fn
	lookupTXT = fn

	return func() { lookupTXT = old }
}
