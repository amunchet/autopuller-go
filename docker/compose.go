package docker

import (
	"context"
	"log"
	"os"
	"os/exec"
)

// DockerManager is an interface for Docker-related operations.
type DockerManager interface {
	RestartServices(ctx context.Context) error
}

type RealDockerManager struct{}

// commandContext is a wrapper around exec.CommandContext, allowing it to be mocked in tests.
var commandContext = exec.CommandContext

// runCommand executes a shell command and logs the output.
func runCommand(ctx context.Context, name string, args ...string) error {

	cmd := commandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Command %s failed: %v", name, err)
		return err
	}
	return nil
}

// RestartServices runs `docker-compose build` and `docker-compose up -d`.
func (d *RealDockerManager) RestartServices(ctx context.Context) error {
	// Change to the directory where the Docker Compose file is located
	repoDir := os.Getenv("DOCKERDIR")
	if err := os.Chdir(repoDir); err != nil {
		return err
	}

	dockercommand := os.Getenv("DOCKERCOMMAND")
	if dockercommand == "" {
		dockercommand = "docker-compose"
	}

	log.Printf("Running %s build...\n", dockercommand)

	// Load in if there's a docker override

	if err := runCommand(ctx, "bash", "-c", dockercommand+" build"); err != nil {
		return err
	}

	log.Printf("Running %s start...\n", dockercommand)
	if err := runCommand(ctx, "bash", "-c", dockercommand+" start"); err != nil {
		return err
	}
	log.Printf("Running %s restart\n", dockercommand)

	if err := runCommand(ctx, "bash", "-c", dockercommand+" restart"); err != nil {
		return err
	}
	return nil
}
