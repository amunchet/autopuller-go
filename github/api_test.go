package github

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Mock HTTP client using httptest to simulate GitHub API responses.
func TestGetMasterSum_Success(t *testing.T) {
	// Set up a fake GitHub API server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/fake-repo/commits/master" {
			t.Errorf("Expected request to '/repos/fake-repo/commits/master', got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"sha": "fake-sha"}`))
	}))
	defer ts.Close()

	// Set environment variables
	os.Setenv("REPONAME", "fake-repo")
	os.Setenv("GITHUBKEY", "fake-key")

	// Override the GitHub API URL with the test server URL
	originalGitHubAPI := fmt.Sprintf("https://api.github.com/repos/%s/commits/master", os.Getenv("REPONAME"))
	overrideGitHubAPI := strings.Replace(originalGitHubAPI, "https://api.github.com", ts.URL, 1)
	os.Setenv("GITHUB_URL_OVERRIDE", overrideGitHubAPI)

	// Create an instance of RealGitHubAPI
	github := &RealGitHubAPI{}

	// Call GetMasterSum and check the result
	sha, err := github.GetMasterSum(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if sha != "fake-sha" {
		t.Fatalf("Expected SHA 'fake-sha', got '%s'", sha)
	}
}

func TestGetMasterSum_Failure(t *testing.T) {
	// Set up a fake GitHub API server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "GitHub API failure", http.StatusInternalServerError)
	}))
	defer ts.Close()

	// Set environment variables
	os.Setenv("REPONAME", "fake-repo")
	os.Setenv("GITHUBKEY", "fake-key")

	// Override the GitHub API URL with the test server URL
	originalGitHubAPI := fmt.Sprintf("https://api.github.com/repos/%s/commits/master", os.Getenv("REPONAME"))
	overrideGitHubAPI := strings.Replace(originalGitHubAPI, "https://api.github.com", ts.URL, 1)
	os.Setenv("GITHUB_URL_OVERRIDE", overrideGitHubAPI)

	// Create an instance of RealGitHubAPI
	github := &RealGitHubAPI{}

	// Call GetMasterSum and expect an error
	_, err := github.GetMasterSum(context.Background())
	if err == nil {
		t.Fatalf("Expected an error, but got nil")
	}
}

// Test GetCurrentSum by mocking file system operations.
func TestGetCurrentSum_Success(t *testing.T) {
	// Manually create a temporary directory
	repoDir, err := ioutil.TempDir("", "repo")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(repoDir) // Clean up after the test

	// Mock REPODIR environment variable
	os.Setenv("REPODIR", repoDir)

	// Create a fake master file with a commit SHA
	masterFile := filepath.Join(repoDir, ".git/refs/heads/master")
	os.MkdirAll(filepath.Dir(masterFile), 0755)
	ioutil.WriteFile(masterFile, []byte("fake-local-sha"), 0644)

	// Create an instance of RealGitHubAPI
	github := &RealGitHubAPI{}

	// Call GetCurrentSum and check the result
	sha, err := github.GetCurrentSum()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if sha != "fake-local-sha" {
		t.Fatalf("Expected SHA 'fake-local-sha', got '%s'", sha)
	}
}

// Test UpdateCurrentSum by mocking file system operations.
func TestUpdateCurrentSum_Success(t *testing.T) {
	// Manually create a temporary directory
	repoDir, err := ioutil.TempDir("", "repo")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(repoDir) // Clean up after the test

	// Mock REPODIR environment variable
	os.Setenv("REPODIR", repoDir)

	// Create a fake master file path
	masterFile := filepath.Join(repoDir, ".git/refs/heads/master")
	os.MkdirAll(filepath.Dir(masterFile), 0755)

	// Create an instance of RealGitHubAPI
	github := &RealGitHubAPI{}

	// Call UpdateCurrentSum to write a new SHA to the master file
	err = github.UpdateCurrentSum("new-fake-sha")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify that the file contains the new SHA
	data, err := ioutil.ReadFile(masterFile)
	if err != nil {
		t.Fatalf("Expected no error reading the file, got %v", err)
	}
	if strings.TrimSpace(string(data)) != "new-fake-sha" {
		t.Fatalf("Expected SHA 'new-fake-sha', got '%s'", strings.TrimSpace(string(data)))
	}
}

// Test CheckLastRun with mocked GitHub Actions API.
func TestCheckLastRun_Success(t *testing.T) {
	// Set up a fake GitHub API server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/fake-repo/actions/runs" {
			t.Errorf("Expected request to '/repos/fake-repo/actions/runs', got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"workflow_runs": [
				{
					"head_sha": "fake-sha",
					"conclusion": "success"
				}
			]
		}`))
	}))
	defer ts.Close()

	// Set environment variables
	os.Setenv("REPONAME", "fake-repo")
	os.Setenv("GITHUBKEY", "fake-key")

	// Override the GitHub API URL with the test server URL
	originalGitHubAPI := fmt.Sprintf("https://api.github.com/repos/%s/actions/runs", os.Getenv("REPONAME"))
	overrideGitHubAPI := strings.Replace(originalGitHubAPI, "https://api.github.com", ts.URL, 1)
	os.Setenv("GITHUB_URL_OVERRIDE", overrideGitHubAPI)

	// Create an instance of RealGitHubAPI
	github := &RealGitHubAPI{}

	// Call CheckLastRun and check the result
	success, err := github.CheckLastRun(context.Background(), "fake-sha")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !success {
		t.Fatalf("Expected success, but got failure")
	}
}