package github

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// Mock HTTP client using httptest to simulate GitHub API responses.

func TestGetMasterSum_Success(t *testing.T) {
	// Set up a fake GitHub API server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/fake-repo/commits/master" {
			t.Errorf("Expected request to '/repos/fake-repo/commits/master', got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"sha": "fake-sha"}`))
	}))
	defer ts.Close()

	// Set environment variables
	os.Setenv("REPONAME", "fake-repo") // Mock the repository name
	os.Setenv("GITHUBKEY", "fake-key") // Mock the GitHub API key

	// Set the environment variable to override the URL prefix with the test server URL
	os.Setenv("GITHUB_URL_PREFIX", ts.URL+"/")

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

	// Set the environment variable to override the URL prefix with the test server URL
	os.Setenv("GITHUB_URL_PREFIX", ts.URL+"/")

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

	// Set the environment variable to override the URL prefix with the test server URL
	os.Setenv("GITHUB_URL_PREFIX", ts.URL+"/repos/")

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

func TestCheckDifferences_Success(t *testing.T) {
	// Set up a fake GitHub API server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request URL
		if r.URL.Path != "/repos/fake-repo/compare/oldSha...newSha" {
			t.Errorf("Expected request to '/repos/fake-repo/compare/oldSha...newSha', got '%s'", r.URL.Path)
		}

		// Respond with a fake list of changed files
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"files": [
				{"filename": "file1.txt"},
				{"filename": "file2.txt"}
			]
		}`))
	}))
	defer ts.Close()

	// Set environment variables for the test
	os.Setenv("REPONAME", "fake-repo")
	os.Setenv("GITHUBKEY", "fake-key")

	// Use the test server's URL as a mock GitHub API endpoint
	os.Setenv("GITHUB_URL_PREFIX", ts.URL+"/repos/")

	// Create an instance of RealGitHubAPI
	github := &RealGitHubAPI{}

	// Call CheckDifferences and check the result
	diffs, err := github.CheckDifferences(context.Background(), "oldSha", "newSha")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if the returned list of differences matches the expected output
	expectedDiffs := []string{"file1.txt", "file2.txt"}
	if len(diffs) != len(expectedDiffs) {
		t.Fatalf("Expected %d differences, got %d", len(expectedDiffs), len(diffs))
	}
	for i, diff := range diffs {
		if diff != expectedDiffs[i] {
			t.Fatalf("Expected difference '%s', got '%s'", expectedDiffs[i], diff)
		}
	}
}
