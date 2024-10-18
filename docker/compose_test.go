package docker

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"testing"
)

// mockExecCommand simulates the behavior of exec.CommandContext for testing.
var mockExecCommand = func(ctx context.Context, name string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", name}
	cs = append(cs, args...)
	cmd := exec.CommandContext(ctx, os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// TestHelperProcess is used as a helper process to simulate exec.CommandContext behavior.
// This simulates a successful execution of the command.
func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// Simulate success by exiting with code 0
	os.Exit(0)
}

// TestHelperProcessFail simulates a failure in exec.CommandContext.
func TestHelperProcessFail(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// Simulate failure by exiting with code 1
	os.Exit(1)
}

// Override exec.CommandContext with the mockExecCommand in tests.
func init() {
	execCommandContext = mockExecCommand
}

// TestRestartServices_Success tests the RestartServices method when the commands succeed.
func TestRestartServices_Success(t *testing.T) {
	// Redirect output to a buffer to capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	// Mock DOCKERDIR environment variable
	os.Setenv("DOCKERDIR", "/fake/dir")

	// Create a new RealDockerManager instance
	dockerMgr := &RealDockerManager{}

	// Call RestartServices and check if it succeeds
	err := dockerMgr.RestartServices(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	// Check if log output contains the expected messages
	expectedLog := "Running docker-compose build...\nRunning docker-compose up -d...\n"
	if !bytes.Contains(buf.Bytes(), []byte(expectedLog)) {
		t.Fatalf("Expected log output to contain '%s', but got '%s'", expectedLog, buf.String())
	}
}

// TestRestartServices_Failure tests the RestartServices method when the commands fail.
func TestRestartServices_Failure(t *testing.T) {
	// Modify the mockExecCommand to simulate a failure by using TestHelperProcessFail
	mockExecCommand = func(ctx context.Context, name string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcessFail", "--", name}
		cs = append(cs, args...)
		cmd := exec.CommandContext(ctx, os.Args[0], cs...)
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
		return cmd
	}

	// Create a new RealDockerManager instance
	dockerMgr := &RealDockerManager{}

	// Call RestartServices and check if it fails
	err := dockerMgr.RestartServices(context.Background())
	if err == nil {
		t.Fatal("Expected an error, but got none")
	}
}
