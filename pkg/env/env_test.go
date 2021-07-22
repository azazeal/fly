package env

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"testing"

	"github.com/azazeal/fly/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	nonFlyKey1 = "SOME_NON_FLY_RANDOM_KEY"
	nonFlyKey2 = "ANOTHER_NON_FLY_RANDOM_KEY"
)

func TestIsSet(t *testing.T) {
	forEachCase(t, func(t *testing.T, kase *testCase) {
		defer testutil.SetEnv(t, kase.env)()

		assert.Equal(t, kase.exp, IsSet())
	})
}

func TestMap(t *testing.T) {
	forEachCase(t, func(t *testing.T, kase *testCase) {
		var exp map[string]string

		for k, v := range kase.env {
			if isNonFlyKey(k) {
				continue
			}
			if exp == nil {
				exp = make(map[string]string)
			}
			exp[k] = v
		}

		defer testutil.SetEnv(t, kase.env)()

		assert.Equal(t, exp, Map())
	})
}

func forEachCase(t *testing.T, fn func(*testing.T, *testCase)) {
	t.Helper()

	cases := buildCases(t)

	for caseIndex := range cases {
		kase := cases[caseIndex]

		t.Run(strconv.Itoa(caseIndex), func(t *testing.T) {
			fn(t, kase)
		})
	}
}

func TestGetters(t *testing.T) {
	funcs := map[string]func() string{
		AppNameKey:  AppName,
		AllocIDKey:  AllocID,
		PublicIPKey: PublicIP,
		RegionKey:   Region,
	}

	for key := range funcs {
		get := funcs[key]

		t.Run(key, func(t *testing.T) {
			defer testutil.SetEnv(t, nil)()
			assert.Equal(t, "", get())

			exp := value(t)
			_ = testutil.SetEnv(t, map[string]string{key: exp})
			assert.Equal(t, exp, get())
		})
	}
}

type testCase struct {
	env map[string]string
	exp bool
}

func buildCases(t *testing.T) []*testCase {
	t.Helper()

	return []*testCase{
		0: {},
		1: {
			env: map[string]string{},
		},
		2: {
			env: map[string]string{
				nonFlyKey1: value(t),
			},
		},
		3: {
			env: map[string]string{
				nonFlyKey1: value(t),
				AppNameKey: value(t),
			},
		},
		4: {
			env: map[string]string{
				nonFlyKey1: value(t),
				AppNameKey: value(t),
				AllocIDKey: value(t),
			},
		},
		5: {
			env: map[string]string{
				nonFlyKey1:  value(t),
				AppNameKey:  value(t),
				AllocIDKey:  value(t),
				PublicIPKey: value(t),
			},
		},
		6: {
			env: map[string]string{
				nonFlyKey1:  value(t),
				AppNameKey:  value(t),
				AllocIDKey:  value(t),
				PublicIPKey: value(t),
				RegionKey:   value(t),
			},
			exp: true,
		},
		7: {
			env: map[string]string{
				nonFlyKey1:  value(t),
				AppNameKey:  value(t),
				AllocIDKey:  value(t),
				PublicIPKey: value(t),
				RegionKey:   value(t),
				nonFlyKey2:  value(t),
			},
			exp: true,
		},
	}
}

func value(t *testing.T) string {
	t.Helper()

	buf := make([]byte, 5)
	_, err := rand.Read(buf)
	require.NoError(t, err)

	return hex.EncodeToString(buf)
}

func isNonFlyKey(key string) bool {
	if key == nonFlyKey1 || key == nonFlyKey2 {
		return true
	}

	return false
}
