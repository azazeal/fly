package testutil

import (
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertHeadersAreCanonical asserts the given named headers are in canonical
// format.
func AssertHeadersAreCanonical(tb testing.TB, headers map[string]string) (ok bool) {
	tb.Helper()

	require.NotEmpty(tb, headers)

	ok = true

	for name, val := range headers {
		exp := textproto.CanonicalMIMEHeaderKey(val)

		if !assert.Equal(tb, exp, val, "header: %s", name) {
			ok = false
		}
	}

	return
}
