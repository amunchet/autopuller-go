package github

import (
	"context"
	"errors"
)

// MockGitHubAPI is a mock implementation of the GitHubAPI interface for testing purposes.
type MockGitHubAPI struct {
	// Fields to override the return values of the methods for specific test scenarios
	OverrideMasterSum    string
	OverrideCurrentSum   string
	OverrideCheckLastRun bool
	ShouldFailMasterSum  bool
	ShouldFailCurrentSum bool
	ShouldFailCheckRun   bool
	ShouldFailUpdateSum  bool
}

// GetMasterSum simulates fetching the latest commit SHA from GitHub.
func (m *MockGitHubAPI) GetMasterSum(ctx context.Context) (string, error) {
	if m.ShouldFailMasterSum {
		return "", errors.New("failed to get master sum")
	}
	return m.OverrideMasterSum, nil
}

// GetCurrentSum simulates reading the current commit SHA from the local file system.
func (m *MockGitHubAPI) GetCurrentSum() (string, error) {
	if m.ShouldFailCurrentSum {
		return "", errors.New("failed to get current sum")
	}
	return m.OverrideCurrentSum, nil
}

// CheckLastRun simulates checking if the last GitHub Actions run succeeded.
func (m *MockGitHubAPI) CheckLastRun(ctx context.Context, sha string) (bool, error) {
	if m.ShouldFailCheckRun {
		return false, errors.New("failed to check last run")
	}
	return m.OverrideCheckLastRun, nil
}

// UpdateCurrentSum simulates updating the current commit SHA in the local file system.
func (m *MockGitHubAPI) UpdateCurrentSum(sha string) error {
	if m.ShouldFailUpdateSum {
		return errors.New("failed to update current sum")
	}
	return nil
}
