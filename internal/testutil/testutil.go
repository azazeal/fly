// Package testutil implements functionality the test suites consume.
package testutil

import (
	"bytes"
	"reflect"
	"testing"
)

// AssertEqual asserts that want & have are equal.
func AssertEqual(tb testing.TB, expected, actual any) (are bool) {
	tb.Helper()

	if are = areEqual(expected, actual); !are {
		tb.Errorf("Not equal: \nexpected: %#v\nactual  : %#v", expected, actual)
	}

	return
}

func areEqual(expected, actual any) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}
