package main

import (
	"os"
	"reflect"
	"testing"

	"hermannm.dev/devlog"
	"hermannm.dev/devlog/log"
)

// This function runs before all test functions, to set up our logger.
func TestMain(m *testing.M) {
	log.SetDefault(devlog.NewHandler(os.Stdout, &devlog.Options{ForceColors: true}))

	os.Exit(m.Run())
}

func assertEqual[T comparable](t *testing.T, actual T, expected T, descriptor string) {
	t.Helper()

	if actual != expected {
		t.Errorf(
			`Unexpected %s
Want: %+v
 Got: %+v`,
			descriptor,
			expected,
			actual,
		)
	}
}

func assertDeepEqual(t *testing.T, actual any, expected any, descriptor string) {
	t.Helper()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf(
			`Unexpected %s
Want: %+v
 Got: %+v`,
			descriptor,
			expected,
			actual,
		)
	}
}

func assertNilError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
