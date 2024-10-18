package github

import (
	"context"
	"os"
	"os/exec"
	"testing"
)

// mockExecCommand is used to replace exec.CommandContext in the test to avoid executing real commands.
var mockExecCommand = func(ctx context.Context, name string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", name}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// TestHelperProcess is used as a helper to mock exec.CommandContext during tests.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	// Simulate successful command execution
	os.Exit(0)
}

// TestRunGitPull_Success simulates the success scenario of the RunGitPull function.
func TestRunGitPull_Success(t *testing.T) {
	// Save original chdir and execCommandContext
	originalChdir := chdir
	originalExecCommandContext := execCommandContext

	// Restore them after the test
	defer func() {
		chdir = originalChdir
		execCommandContext = originalExecCommandContext
	}()

	// Mock chdir to always succeed
	chdir = func(dir string) error {
		return nil
	}

	// Mock execCommandContext to use mockExecCommand instead of running real commands
	execCommandContext = mockExecCommand

	// Create an instance of RealGitHubAPI and call RunGitPull
	github := &RealGitHubAPI{}
	err := github.RunGitPull(context.Background(), "/path/to/repo")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

// TestRunGitPull_Failure simulates a failure scenario where a git command fails.
func TestRunGitPull_Failure(t *testing.T) {
	// Save original chdir and execCommandContext
	originalChdir := chdir
	originalExecCommandContext := execCommandContext

	// Restore them after the test
	defer func() {
		chdir = originalChdir
		execCommandContext = originalExecCommandContext
	}()

	// Mock chdir to always succeed
	chdir = func(dir string) error {
		return nil
	}

	// Mock execCommandContext to simulate a failure for one of the git commands
	execCommandContext = func(ctx context.Context, name string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcessFail", "--", name}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS_FAIL=1"}
		return cmd
	}

	// Create an instance of RealGitHubAPI and call RunGitPull
	github := &RealGitHubAPI{}
	err := github.RunGitPull(context.Background(), "/path/to/repo")

	// We expect an error here because of the simulated failure
	if err == nil {
		t.Fatalf("Expected error, but got nil")
	}
}

// TestHelperProcessFail simulates a command failure scenario.
func TestHelperProcessFail(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS_FAIL") != "1" {
		return
	}

	// Simulate a failure by returning non-zero exit code
	os.Exit(1)
}
