package util

import (
	"os"
	"testing"
)

func IntegrationTest(t *testing.T) {
	t.Helper()
	dir, err := os.MkdirTemp("", "nvmc-home")
	if err != nil {
		t.Fatalf("Failed to make temp NVMC_HOME directory: %v", err)
	}
	defer os.RemoveAll(dir)
	if err := os.Setenv("NVMC_HOME", dir); err != nil {
		t.Fatalf("Failed to set NVMC_HOME: %v", err)
	}
}
