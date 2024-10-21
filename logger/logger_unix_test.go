//go:build !windows
// +build !windows

package logger

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

// TestInitLogger_Unix tests logger initialization on Unix-like systems.
func TestInitLogger_Unix(t *testing.T) {
	// Redirect log output to a temporary file to capture it
	tempFile, err := ioutil.TempFile("", "log_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up

	// Override log output to the temporary file
	log.SetOutput(tempFile)

	// Initialize the logger (this includes syslog and file logging)
	InitLogger("application.log")

	// Check if the log file contains the initialization message
	content, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	if !contains(content, "Logger initialized") {
		t.Fatalf("Expected 'Logger initialized' in log, but it was not found")
	}

	// Optionally, check for syslog behavior, but it's generally not easily testable
	// in a unit test environment.
}

// Helper function to check if a string is contained in a byte slice.
func contains(content []byte, substr string) bool {
	return string(content) == substr
}
