package testutil

import (
	"io"
	"testing"
)

// ReadAll reads r into a string.
//
// Should an error occur while reading from r, ReadAll will fail tb.
func ReadAll(tb testing.TB, r io.Reader) string {
	tb.Helper()

	data, err := io.ReadAll(r)
	if err != nil {
		tb.Fatalf("failed reading: %v", err)
	}
	return string(data)
}
