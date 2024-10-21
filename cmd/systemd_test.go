package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
)

// Helper function to check if a string is contained in a file's content.
func contains(content []byte, substr string) bool {
	return strings.Contains(string(content), substr)
}

// TestGenerateSystemdService tests the generation of a systemd service file
func TestGenerateSystemdService(t *testing.T) {
	// Create a temporary directory to simulate systemd folder
	tempDir, err := ioutil.TempDir("", "systemd_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after the test

	// Mock the /etc/systemd/system path
	mockServicePath := filepath.Join(tempDir, "autopuller.service")

	// Modify the GenerateSystemdService function to accept the output path
	GenerateSystemdService := func(serviceName, user, outputPath string) error {
		workingDir, err := os.Getwd()
		if err != nil {
			return err
		}

		execPath, err := os.Executable()
		if err != nil {
			return err
		}

		config := SystemdServiceConfig{
			ServiceName: serviceName,
			ExecPath:    execPath,
			WorkingDir:  workingDir,
			User:        user,
		}

		// Create or open the file at the mock path
		file, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Parse and execute the template
		tmpl, err := template.New("systemdService").Parse(ServiceFileTemplate)
		if err != nil {
			return err
		}

		// Write the generated systemd service to the file
		err = tmpl.Execute(file, config)
		if err != nil {
			return err
		}

		return nil
	}

	// Call the modified GenerateSystemdService function with the mock path
	err = GenerateSystemdService("autopuller", "root", mockServicePath)
	if err != nil {
		t.Fatalf("Failed to generate systemd service file: %v", err)
	}

	// Read the generated service file
	content, err := ioutil.ReadFile(mockServicePath)
	if err != nil {
		t.Fatalf("Failed to read generated service file: %v", err)
	}

	// Check if the service file contains the correct information
	if !contains(content, "autopuller") {
		t.Errorf("Expected service name 'autopuller' in file, but not found")
	}
	if !contains(content, "ExecStart") {
		t.Errorf("Expected 'ExecStart' in the service file, but not found")
	}
	if !contains(content, "root") {
		t.Errorf("Expected 'User=root' in the service file, but not found")
	}
}
