package start

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/docker/docker/api/types"
)

func TestStartSuccess(t *testing.T) {
	// Arrange
	containers := []types.Container{
		{
			ID:    "nitro",
			Names: []string{"/nitro"},
			//Labels: map[string]string{
			//	containerlabels.Host: "nitro",
			//},
		},
	}
	expectedContainerID := "nitro"
	mock := newMockDockerClient(nil, containers, nil)
	output := &spyOutputer{}
	expectedOutput := []string{"Starting Nitro…\n", "Nitro started 👍\n"}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	home := filepath.Join(wd, "testdata")

	// Act
	cmd := NewCommand(home, mock, output)
	err = cmd.RunE(cmd, []string{})

	// Assert
	if err != nil {
		t.Errorf("expected the error to be nil, got %v", err)
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
	expectedOutputSuccess := []string{"  ✓ testing-start\n"}
	home, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Act
	cmd := NewCommand(home, mock, output)
	err = cmd.RunE(cmd, []string{})

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
	mock := newMockDockerClient(nil, nil, nil)
	home, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Act
	cmd := NewCommand(home, mock, &spyOutputer{})
	err = cmd.RunE(cmd, os.Args)

	// Assert
	if err == nil {
		t.Errorf("expected the error to not be nil")
	}
}
