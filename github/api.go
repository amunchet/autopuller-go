package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GitHubAPI is an interface that defines the functions interacting with GitHub.
// This is for testing
type GitHubAPI interface {
	GetMasterSum(ctx context.Context) (string, error)
	GetCurrentSum() (string, error)
	CheckLastRun(ctx context.Context, sha string) (bool, error)
	CheckDifferences(ctx context.Context, oldSha, newSha string) ([]string, error)
	RunGitPull(ctx context.Context, repoDir string) error
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
	url := fmt.Sprintf("%s%s/commits/master", urlPrefix, repoName)

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
	// Define default GitHub API URL prefix
	defaultURLPrefix := "https://api.github.com/repos/"

	// Check if we have an override for the URL prefix (for testing purposes)
	urlPrefix := os.Getenv("GITHUB_URL_PREFIX")
	if urlPrefix == "" {
		urlPrefix = defaultURLPrefix
	}
	// Construct the final URL for the GitHub Actions API
	url := fmt.Sprintf("%s%s/actions/runs", urlPrefix, os.Getenv("REPONAME"))

	log.Println("URL:", url)
	// Create a new HTTP GET request
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUBKEY"))

	// Send the HTTP request and handle errors
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal("HTTP Error:", err)
		return false, err
	}
	defer resp.Body.Close()

	// Check for HTTP OK
	if resp.StatusCode != 200 {
		log.Println("Unexpected status code: ", resp.StatusCode)
		return false, err
	}

	// Parse the response body
	var result struct {
		WorkflowRuns []struct {
			HeadSha    string `json:"head_sha"`
			Conclusion string `json:"conclusion"`
		} `json:"workflow_runs"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal("Error:", err)
		return false, err
	}

	// Check if the run corresponding to the given SHA was successful

	for _, run := range result.WorkflowRuns {

		log.Println("Run:", run)
		if run.HeadSha == sha && run.Conclusion == "success" {
			return true, nil
		}
	}

	return false, nil
}

// CheckDifferences compares two SHAs and returns a list of changed files.
func (g *RealGitHubAPI) CheckDifferences(ctx context.Context, oldSha, newSha string) ([]string, error) {
	// Use GITHUB_URL_PREFIX if set, otherwise default to the GitHub API URL.
	urlPrefix := os.Getenv("GITHUB_URL_PREFIX")
	if urlPrefix == "" {
		urlPrefix = "https://api.github.com/repos/"
	}

	// Construct the full URL using the prefix and the SHAs.
	url := fmt.Sprintf("%s%s/compare/%s...%s", urlPrefix, os.Getenv("REPONAME"), oldSha, newSha)

	log.Println("CheckDifferences URL:", url)

	// Create the HTTP request
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUBKEY"))

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for HTTP OK
	if resp.StatusCode != 200 {
		log.Println("Unexpected status code: ", resp.StatusCode)
		return nil, err
	}

	// Decode the JSON response
	var result struct {
		Files []struct {
			Filename string `json:"filename"`
		} `json:"files"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Collect the filenames of the changed files
	var diffs []string
	for _, file := range result.Files {
		diffs = append(diffs, file.Filename)
	}

	return diffs, nil
}

// Define function variables that can be overridden in tests
var chdir = os.Chdir
var execCommandContext = exec.CommandContext

// RunGitPull runs git-related commands to update the repository.
func (g *RealGitHubAPI) RunGitPull(ctx context.Context, repoDir string) error {
	// Change directory to the repoDir
	if err := chdir(repoDir); err != nil {
		return err
	}

	// Set git credential helper and safe directory
	commands := [][]string{
		{"git", "config", "credential.helper", "store"},
		{"git", "config", "--global", "--add", "safe.directory", repoDir},
		{"git", "pull"},
	}

	for _, cmdArgs := range commands {
		cmd := execCommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error running command %s: %v", cmdArgs[0], err)
			log.Printf("Output: %s", output)
			return fmt.Errorf("failed to run git command: %s", cmdArgs)
		}
		log.Printf("Command %s successful: %s", cmdArgs[0], string(output))
	}

	return nil
}
