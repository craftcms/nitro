package client

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
)

func TestStopSuccess(t *testing.T) {
	// Arrange
	environmentName := "testing-stop"
	containers := []types.Container{
		{
			ID:    "testing-stop",
			Names: []string{"/testing-stop"},
		},
	}
	expectedContainerID := "testing-stop"
	mock := newMockDockerClient(nil, containers, nil)
	cli := Client{docker: mock}
	ctx := context.TODO()

	// Act
	if err := cli.Stop(ctx, environmentName, []string{}); err != nil {
		t.Errorf("expected the error to be nil, got %w", err)
	}

	// Assert
	if mock.containerID != expectedContainerID {
		t.Errorf("expected the container IDs to match got %s, want %s", mock.containerID, expectedContainerID)
	}
}

func TestStopErrorsWhenThereAreNoContainers(t *testing.T) {
	// Arrange
	environmentName := "testing-stop"
	mock := newMockDockerClient(nil, nil, nil)
	cli := Client{docker: mock}
	ctx := context.TODO()

	// Act
	err := cli.Stop(ctx, environmentName, nil)
	if err != nil {
		t.Errorf("expected the error to not be nil")
	}
}
