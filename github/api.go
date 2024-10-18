package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// GitHubAPI is an interface that defines the functions interacting with GitHub.
// This is for testing
type GitHubAPI interface {
	GetMasterSum(ctx context.Context) (string, error)
	GetCurrentSum() (string, error)
	CheckLastRun(ctx context.Context, sha string) (bool, error)
	UpdateCurrentSum(sha string) error
}

type RealGitHubAPI struct {
}

// GetMasterSum fetches the latest commit SHA from GitHub for the master branch.
func (g *RealGitHubAPI) GetMasterSum(ctx context.Context) (string, error) {
	// Define default URL prefix
	defaultURLPrefix := "https://api.github.com/repos/"

	// Check if we have a URL override for the prefix, otherwise use the default
	urlPrefix := os.Getenv("GITHUB_URL_PREFIX")
	if urlPrefix == "" {
		urlPrefix = defaultURLPrefix
	}

	// Retrieve the repository name from environment variables
	repoName := os.Getenv("REPONAME")
	if repoName == "" {
		return "", fmt.Errorf("REPONAME environment variable not set")
	}

	// Construct the final URL using the prefix and repository name
	url := fmt.Sprintf("%s/%s/commits/master", urlPrefix, repoName)

	log.Println("Request URL:", url)

	// Prepare the request
	req, _ := http.NewRequest("GET", url, nil)
	if githubkey := os.Getenv("GITHUBKEY"); githubkey != "" {
		req.Header.Set("Authorization", "token "+githubkey)
	}

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode the response
	var result struct {
		Sha string `json:"sha"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Sha, nil
}

// GetCurrentSum reads the current commit SHA from the local file system.
func (g *RealGitHubAPI) GetCurrentSum() (string, error) {
	repoDir := os.Getenv("REPODIR")
	if err := os.Chdir(repoDir); err != nil {
		return "", err
	}

	var masterFile = filepath.Join(os.Getenv("REPODIR"), ".git/refs/heads/master")
	filename := filepath.FromSlash(masterFile)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// CheckLastRun checks if the last GitHub Actions run for the commit was successful.
func (g *RealGitHubAPI) CheckLastRun(ctx context.Context, sha string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/actions/runs", os.Getenv("REPONAME"))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUBKEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result struct {
		WorkflowRuns []struct {
			HeadSha    string `json:"head_sha"`
			Conclusion string `json:"conclusion"`
		} `json:"workflow_runs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	for _, run := range result.WorkflowRuns {
		if run.HeadSha == sha && run.Conclusion == "success" {
			return true, nil
		}
	}

	return false, nil
}

// UpdateCurrentSum writes the latest commit SHA to the local file system.
func (g *RealGitHubAPI) UpdateCurrentSum(sha string) error {

	var masterFile = filepath.Join(os.Getenv("REPODIR"), ".git/refs/heads/master")
	log.Println(masterFile)
	return ioutil.WriteFile(masterFile, []byte(sha), 0644)
}
