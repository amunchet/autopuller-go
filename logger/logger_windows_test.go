//go:build windows
// +build windows

package logger

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

// Helper function to check if a string is contained in a byte slice.
func contains(content []byte, substr string) bool {
	// Convert content to string and check if substr is a part of it
	return strings.Contains(string(content), substr)
}

// TestInitLogger_Windows tests logger initialization on Windows.
func TestInitLogger_Windows(t *testing.T) {
	// Redirect log output to a temporary file to capture it
	tempFile, err := ioutil.TempFile("", "log_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up

	// Initialize the logger (this excludes syslog)
	InitLogger(tempFile.Name())

	// Check if the log file contains the initialization message
	content, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}
	log.Println(content)
	if !contains(content, "Logger initialized") {
		t.Fatalf("Expected 'Logger initialized' in log, but it was not found")
	}
}

// Helper function to check if a string is contained in a byte slice
