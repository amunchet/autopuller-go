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

const masterFile = ".git/refs/heads/master"

// GetMasterSum fetches the latest commit SHA from GitHub for the master branch.
func GetMasterSum(ctx context.Context) (string, error) {

	url := fmt.Sprintf("https://api.github.com/repos/%s/commits/master", os.Getenv("REPONAME"))

	log.Println(url)

	req, _ := http.NewRequest("GET", url, nil)
	if githubkey := os.Getenv("GITHUBKEY"); githubkey != "" {
		req.Header.Set("Authorization", "token "+githubkey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Sha string `json:"sha"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Sha, nil
}

// GetCurrentSum reads the current commit SHA from the local file system.
func GetCurrentSum() (string, error) {
	repoDir := os.Getenv("REPODIR")
	if err := os.Chdir(repoDir); err != nil {
		return "", err
	}

	filename := filepath.FromSlash(masterFile)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// CheckLastRun checks if the last GitHub Actions run for the commit was successful.
func CheckLastRun(ctx context.Context, sha string) (bool, error) {
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
func UpdateCurrentSum(sha string) error {
	return ioutil.WriteFile(masterFile, []byte(sha), 0644)
}
