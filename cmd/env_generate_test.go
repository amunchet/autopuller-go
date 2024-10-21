package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// TestGenerateEnvSample tests the generation of the .env.sample file
func TestGenerateEnvSample(t *testing.T) {
	// Create a temporary directory to simulate the current working directory
	tempDir, err := ioutil.TempDir("", "envsample_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after the test

	// Change the working directory to the temporary directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(oldDir) // Restore the original working directory after the test

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}

	// Call the GenerateEnvSample function
	err = GenerateEnvSample()
	if err != nil {
		t.Fatalf("Failed to generate .env.sample file: %v", err)
	}

	// Check if the .env.sample file was created
	envSamplePath := filepath.Join(tempDir, ".env.sample")
	_, err = os.Stat(envSamplePath)
	if os.IsNotExist(err) {
		t.Fatalf(".env.sample file was not created")
	}

	// Read the contents of the .env.sample file
	content, err := ioutil.ReadFile(envSamplePath)
	if err != nil {
		t.Fatalf("Failed to read .env.sample file: %v", err)
	}

	// Check if the file contains expected content
	if !contains(content, "GITHUBKEY=") {
		t.Errorf("Expected GITHUBKEY field in .env.sample, but not found")
	}
	if !contains(content, "REPONAME=amunchet/autopuller-go") {
		t.Errorf("Expected REPONAME field in .env.sample, but not found")
	}
	if !contains(content, "DOCKERDIR=./docker/sample") {
		t.Errorf("Expected DOCKERDIR field in .env.sample, but not found")
	}
	if !contains(content, "FORCEPULL=") {
		t.Errorf("Expected FORCEPULL field in .env.sample, but not found")
	}
}
