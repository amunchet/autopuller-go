package main

import (
	"context"
	"testing"

	"autopuller/docker"
	"autopuller/github"
)

// TestCheckForUpdates_Success tests the happy path scenario where everything works as expected:
// - A new commit is detected.
// - The last GitHub action run was successful.
// - Docker services restart successfully.
func TestCheckForUpdates_Success(t *testing.T) {
	// Mock GitHub API returning a new master commit
	mockGitHub := &github.MockGitHubAPI{
		// Customize mock responses as needed
	}

	// Mock DockerManager simulating a successful service restart
	mockDocker := &docker.MockDockerManager{
		ShouldFail: false, // Simulate successful restart
	}

	// Call the function under test
	ctx := context.Background()
	err := checkForUpdates(ctx, mockGitHub, mockDocker)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
}

// TestCheckForUpdates_NoNewCommits tests the scenario where there are no new commits to pull from GitHub:
// - The current and master commit SHAs are the same.
// - No services should restart.
func TestCheckForUpdates_NoNewCommits(t *testing.T) {
	// Customize MockGitHubAPI to simulate no new commits
	mockGitHub := &github.MockGitHubAPI{
		OverrideMasterSum:  "same_sha", // Simulate master commit is the same
		OverrideCurrentSum: "same_sha", // Simulate current commit is the same
	}

	// Mock DockerManager should not be called in this case
	mockDocker := &docker.MockDockerManager{
		ShouldFail: false,
	}

	// Call the function under test
	ctx := context.Background()
	err := checkForUpdates(ctx, mockGitHub, mockDocker)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
}

// TestCheckForUpdates_LastRunFailed tests the scenario where the last GitHub action run failed:
// - A new commit is detected.
// - The last GitHub action run failed.
// - Services should not restart.
func TestCheckForUpdates_LastRunFailed(t *testing.T) {
	// Mock GitHub API returning a new commit but with a failed GitHub Actions run
	mockGitHub := &github.MockGitHubAPI{
		OverrideMasterSum:    "new_sha",
		OverrideCurrentSum:   "old_sha",
		OverrideCheckLastRun: false, // Simulate that the last run failed
	}

	// Mock DockerManager should not restart services in this case
	mockDocker := &docker.MockDockerManager{}

	// Call the function under test
	ctx := context.Background()
	err := checkForUpdates(ctx, mockGitHub, mockDocker)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
}

// TestCheckForUpdates_DockerRestartFailure tests the scenario where Docker services fail to restart:
// - A new commit is detected.
// - The last GitHub action run was successful.
// - Docker services fail to restart.
func TestCheckForUpdates_DockerRestartFailure(t *testing.T) {
	// Mock GitHub API returning a new master commit and successful last run
	mockGitHub := &github.MockGitHubAPI{
		OverrideMasterSum:    "new_sha",
		OverrideCurrentSum:   "old_sha",
		OverrideCheckLastRun: true, // Simulate that the last run succeeded
	}

	// Mock DockerManager simulating a service restart failure
	mockDocker := &docker.MockDockerManager{
		ShouldFail: true, // Simulate Docker restart failure
	}

	// Call the function under test
	ctx := context.Background()
	err := checkForUpdates(ctx, mockGitHub, mockDocker)
	if err == nil {
		t.Fatalf("Expected error due to Docker restart failure, but got nil")
	}
}
