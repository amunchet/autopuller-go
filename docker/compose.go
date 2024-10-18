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

// runCommand executes a shell command and logs the output.
func runCommand(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
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

	log.Println("Running docker-compose build...")
	if err := runCommand(ctx, "docker-compose", "build"); err != nil {
		return err
	}

	log.Println("Running docker-compose up -d...")
	return runCommand(ctx, "docker-compose", "up", "-d")
}
