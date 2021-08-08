package testutil

import (
	"os"
	"testing"
)

// SetEnv sets the given environment variables via os.Setenv.
//
// The returned function may be called to revert any changes the call made to
// the application's environment.
func SetEnv(tb testing.TB, env map[string]string) (restore func()) {
	tb.Helper()

	// empty restore func
	restore = func() {}

	for key, newValue := range env {
		oldValue, existed := os.LookupEnv(key)

		if err := os.Setenv(key, newValue); err != nil {
			tb.Fatalf("couldn't set $%s: %v", key, err)
		}

		if existed {
			restore = newEnvSetFunc(tb, key, oldValue, restore)
		} else {
			restore = newEnvUnsetFunc(tb, key, restore)
		}
	}

	return
}

func newEnvSetFunc(tb testing.TB, k, v string, r func()) func() {
	return func() {
		if err := os.Setenv(k, v); err != nil {
			tb.Fatalf("couldn't restore $%s: %v", k, v)
		}

		r()
	}
}

func newEnvUnsetFunc(tb testing.TB, k string, r func()) func() {
	return func() {
		if err := os.Unsetenv(k); err != nil {
			tb.Fatalf("couldn't unset $%s: %v", k, err)
		}

		r()
	}
}
