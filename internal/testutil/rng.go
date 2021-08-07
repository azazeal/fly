package testutil

import (
	"crypto/rand"
	"encoding/binary"
	mand "math/rand"
	"strings"
	"testing"
)

// RNG returns a reference to a Rand that's been seeded with a random seed.
func RNG(t *testing.T) *mand.Rand {
	t.Helper()

	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		t.Fatalf("failed seeding RNG: %v", err)
	}

	seed := int64(binary.BigEndian.Uint64(buf))
	return mand.New(mand.NewSource(seed))
}

// HexString returns a string of l hex characters.
func HexString(t *testing.T, l int) string {
	t.Helper()

	var buf strings.Builder
	buf.Grow(l)

	rng := RNG(t)

	const allowed = "0123456789abcdef"
	for i := 0; i < l; i++ {
		c := allowed[rng.Intn(len(allowed))]
		_ = buf.WriteByte(c)
	}

	return buf.String()
}
