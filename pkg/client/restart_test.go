package client

import (
	"context"
	"reflect"
	"testing"

	"github.com/docker/docker/api/types"
)

func TestRestart(t *testing.T) {
	// Arrange
	environmentName := "testing-restart"
	mock := newMockDockerClient(nil, nil, nil)
	mock.containers = []types.Container{
		{
			ID:    "testing-restart",
			Names: []string{"/testing-restart"},
			Labels: map[string]string{
				EnvironmentLabel:           "testing-restart",
				"com.craftcms.nitro.proxy": "testing-restart",
			},
		},
		{
			ID:    "testing-restart-hostname",
			Names: []string{"/testing-restart-hostname"},
			Labels: map[string]string{
				EnvironmentLabel:           "testing-restart",
				"com.craftcms.nitro.proxy": "testing-restart",
			},
		},
	}
	cli := Client{docker: mock, out: mockOutput{}}

	// Expected
	ids := []string{"testing-restart", "testing-restart-hostname"}

	// Act
	if err := cli.Restart(context.Background(), environmentName, []string{}); err != nil {
		t.Error(err)
	}

	// Assert
	if !reflect.DeepEqual(mock.containerRestartRequests, ids) {
		t.Errorf(
			"expected container restart requests to match\ngot:\n%v\nwant:\n%v",
			mock.containerRestartRequests,
			ids,
		)
	}
}

func TestRestartWithNoContainersDoesNoWork(t *testing.T) {
	// Arrange
	environmentName := "testing-restart"
	mock := newMockDockerClient(nil, nil, nil)
	cli := Client{docker: mock, out: mockOutput{}}

	// Act
	if err := cli.Restart(context.Background(), environmentName, []string{}); err != nil {
		t.Error(err)
	}

	// Assert
	if len(mock.containerRestartRequests) != 0 {
		t.Errorf("expected the number of restart requests to be zero, got %d instead", len(mock.containerRestartRequests))
	}
}
