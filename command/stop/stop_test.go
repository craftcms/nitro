package stop

import (
	"os"
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

	// Act
	cmd := New(mock, spyOutputer{})
	cmd.Flags().StringP("environment", "e", environmentName, "test flag")
	err := cmd.RunE(cmd, os.Args)
	if err != nil {
		t.Errorf("expected the error to be nil")
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

	// Act
	cmd := New(mock, spyOutputer{})
	cmd.Flags().StringP("environment", "e", environmentName, "test flag")
	err := cmd.RunE(cmd, os.Args)
	if err == nil {
		t.Errorf("expected the error to not be nil")
	}
}
