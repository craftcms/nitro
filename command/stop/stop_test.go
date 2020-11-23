package stop

import (
	"os"
	"reflect"
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
	output := &spyOutputer{}
	expectedOutput := []string{"testing-stop shutdown 😴\n"}

	// Act
	cmd := New(mock, output)
	cmd.Flags().StringP("environment", "e", environmentName, "test flag")
	err := cmd.RunE(cmd, os.Args)
	if err != nil {
		t.Errorf("expected the error to be nil")
	}

	// Assert
	if mock.containerID != expectedContainerID {
		t.Errorf("expected the container IDs to match got %s, want %s", mock.containerID, expectedContainerID)
	}
	if !reflect.DeepEqual(output.infos, expectedOutput) {
		t.Errorf("expected the infos to match, got \n%v\nwant:\n%v", output.infos, expectedOutput)
	}
}

func TestStopErrorsWhenThereAreNoContainers(t *testing.T) {
	// Arrange
	environmentName := "testing-stop"
	mock := newMockDockerClient(nil, nil, nil)
	output := &spyOutputer{}

	// Act
	cmd := New(mock, output)
	cmd.Flags().StringP("environment", "e", environmentName, "test flag")
	err := cmd.RunE(cmd, os.Args)

	// Assert
	if err == nil {
		t.Errorf("expected the error to not be nil")
	}
}
