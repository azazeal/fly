package testutil

import (
	"crypto/rand"
	"encoding/binary"
	mand "math/rand"
	"strings"
	"testing"
)

// RNG returns a reference to a Rand that's been seeded with a random seed.
func RNG(tb testing.TB) *mand.Rand {
	tb.Helper()

	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		tb.Fatalf("failed seeding RNG: %v", err)
	}

	seed := int64(binary.BigEndian.Uint64(buf))
	return mand.New(mand.NewSource(seed))
}

// HexString returns a string of l hex characters.
func HexString(tb testing.TB, l int) string {
	tb.Helper()

	var buf strings.Builder
	buf.Grow(l)

	rng := RNG(tb)

	const allowed = "0123456789abcdef"
	for i := 0; i < l; i++ {
		c := allowed[rng.Intn(len(allowed))]
		_ = buf.WriteByte(c)
	}

	return buf.String()
}
