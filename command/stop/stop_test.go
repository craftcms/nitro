package stop

import (
	"os"
	"path/filepath"
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
	expectedOutput := []string{"Stopping Nitro…\n", "Nitro shutdown 😴\n"}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	home := filepath.Join(wd, "testdata")

	// Act
	cmd := NewCommand(home, mock, output)
	cmd.Flags().String("environment", environmentName, "test flag")
	err = cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("expected the error to be nil, got %v", err)
	}

	// Assert
	if mock.containerID != expectedContainerID {
		t.Errorf("expected the container IDs to match got %s, want %s", mock.containerID, expectedContainerID)
	}
	if !reflect.DeepEqual(output.infos, expectedOutput) {
		t.Errorf("expected the infos to match, got \n%v\nwant:\n%v", output.infos, expectedOutput)
	}
}

func TestStopDoesNotErrorWhenThereAreNoContainers(t *testing.T) {
	// Arrange
	environmentName := "testing-stop"
	mock := newMockDockerClient(nil, nil, nil)
	output := &spyOutputer{}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	home := filepath.Join(wd, "testdata")

	// Act
	cmd := NewCommand(home, mock, output)
	cmd.Flags().String("environment", environmentName, "test flag")
	err = cmd.RunE(cmd, []string{})

	// Assert
	if err != nil {
		t.Errorf("expected the error to be nil, got %v", err)
	}
}
