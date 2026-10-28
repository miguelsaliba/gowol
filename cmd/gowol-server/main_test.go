package main

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	const key = "GOWOL_TEST_ENV_VAR"

	t.Run("returns fallback when unset", func(t *testing.T) {
		t.Setenv(key, "") // restores the original value after the test
		os.Unsetenv(key)
		if got := getEnv(key, "fallback"); got != "fallback" {
			t.Errorf("expected %q, got %q", "fallback", got)
		}
	})

	t.Run("returns value when set", func(t *testing.T) {
		t.Setenv(key, "value")
		if got := getEnv(key, "fallback"); got != "value" {
			t.Errorf("expected %q, got %q", "value", got)
		}
	})
}
