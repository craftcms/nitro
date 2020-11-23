package start

import (
	"os"
	"reflect"
	"testing"

	"github.com/docker/docker/api/types"
)

func TestStartSuccess(t *testing.T) {
	// Arrange
	environmentName := "testing-start"
	containers := []types.Container{
		{
			ID:    "testing-start",
			Names: []string{"/testing-start"},
		},
	}
	expectedContainerID := "testing-start"
	mock := newMockDockerClient(nil, containers, nil)
	output := &spyOutputer{}
	expectedOutput := []string{"testing-start started 👍\n"}

	// Act
	cmd := New(mock, output)
	cmd.Flags().StringP("environment", "e", environmentName, "test flag")
	err := cmd.RunE(cmd, os.Args)

	// Assert
	if err != nil {
		t.Errorf("expected the error to be nil")
	}

	if mock.containerID != expectedContainerID {
		t.Errorf("expected the container IDs to match got %s, want %s", mock.containerID, expectedContainerID)
	}

	if !reflect.DeepEqual(output.infos, expectedOutput) {
		t.Errorf("expected the infos to match, got \n%v\nwant:\n%v", output.infos, expectedOutput)
	}
}

func TestStartReturnsReadyIfAlreadyRunning(t *testing.T) {
	// Arrange
	environmentName := "testing-start"
	containers := []types.Container{
		{
			ID:    "testing-start",
			Names: []string{"/testing-start"},
			State: "running",
		},
	}
	expectedContainerID := ""
	mock := newMockDockerClient(nil, containers, nil)
	output := &spyOutputer{}
	expectedOutputSuccess := []string{"  ✓ testing-start ready\n"}

	// Act
	cmd := New(mock, output)
	cmd.Flags().StringP("environment", "e", environmentName, "test flag")
	err := cmd.RunE(cmd, os.Args)

	// Assert
	if err != nil {
		t.Errorf("expected the error to be nil")
	}

	if mock.containerID != expectedContainerID {
		t.Errorf("expected the container IDs to match got %s, want %s", mock.containerID, expectedContainerID)
	}

	if !reflect.DeepEqual(output.succesess, expectedOutputSuccess) {
		t.Errorf("expected the infos to match, got \n%v\nwant:\n%v", output.succesess, expectedOutputSuccess)
	}
}

func TestStartErrorsWhenThereAreNoContainers(t *testing.T) {
	// Arrange
	environmentName := "testing-start"
	mock := newMockDockerClient(nil, nil, nil)

	// Act
	cmd := New(mock, &spyOutputer{})
	cmd.Flags().StringP("environment", "e", environmentName, "test flag")
	err := cmd.RunE(cmd, os.Args)

	// Assert
	if err == nil {
		t.Errorf("expected the error to not be nil")
	}
}
