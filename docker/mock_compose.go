package docker

import (
	"context"
	"errors"
)

// MockDockerManager is a mock implementation of the DockerManager interface for testing purposes.
type MockDockerManager struct {
	// ShouldFail allows us to control whether the mock should simulate a failure when restarting services.
	ShouldFail bool
}

// RestartServices simulates restarting Docker services.
func (m *MockDockerManager) RestartServices(ctx context.Context) error {
	if m.ShouldFail {
		// Simulate a failure
		return errors.New("failed to restart services")
	}
	// Simulate a successful service restart
	return nil
}
