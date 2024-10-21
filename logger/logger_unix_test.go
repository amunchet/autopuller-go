//go:build !windows
// +build !windows

package logger

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// Helper function to check if a string is contained in a byte slice.
func contains(content []byte, substr string) bool {
	return strings.Contains(string(content), substr)
}

// TestInitLogger_Unix tests the logger initialization on Unix-like systems.
func TestInitLogger_Unix(t *testing.T) {
	// Redirect log output to a temporary file to capture it
	tempFile, err := ioutil.TempFile("", "log_test_unix")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name()) // Clean up
	}()

	// Initialize the logger
	err = InitLogger(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Read the log content from the file
	content, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	// Check if the log contains the initialization message
	if !contains(content, "Logger initialized (Unix)") {
		t.Fatalf("Expected 'Logger initialized (Unix)' in log, but it was not found")
	}
}
