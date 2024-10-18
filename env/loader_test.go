package env

import (
	"os"
	"testing"
)

func TestGetInterval(t *testing.T) {
	// Set INTERVAL and check if it's correctly parsed
	os.Setenv("INTERVAL", "30")
	defer os.Unsetenv("INTERVAL")

	interval := GetInterval()
	if interval != 30 {
		t.Fatalf("Expected interval 30, but got %d", interval)
	}

	// Unset INTERVAL and check if it defaults to 60
	os.Unsetenv("INTERVAL")
	interval = GetInterval()
	if interval != 60 {
		t.Fatalf("Expected default interval 60, but got %d", interval)
	}

	// Set invalid INTERVAL and check if it defaults to 60
	os.Setenv("INTERVAL", "invalid")
	interval = GetInterval()
	if interval != 60 {
		t.Fatalf("Expected default interval 60 for invalid INTERVAL, but got %d", interval)
	}
}
