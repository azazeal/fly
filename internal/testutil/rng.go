package testutil

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	mand "math/rand"
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

// Fill fills the given buffer with pseudo-random data.
func Fill(t *testing.T, buf []byte) {
	t.Helper()

	if _, err := io.ReadFull(RNG(t), buf); err != nil {
		t.Fatalf("failed filling buffer: %v", err)
	}
}
