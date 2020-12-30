package restart

import (
	"os"
	"reflect"
	"testing"

	"github.com/docker/docker/api/types"

	"github.com/craftcms/nitro/pkg/labels"
)

func TestRestart(t *testing.T) {
	// Arrange
	mock := newMockDockerClient(nil, nil, nil)
	mock.containers = []types.Container{
		{
			ID:    "testing-restart",
			Names: []string{"/testing-restart"},
			Labels: map[string]string{
				labels.Proxy: "testing-restart",
			},
		},
		{
			ID:    "testing-restart-hostname",
			Names: []string{"/testing-restart-hostname"},
			Labels: map[string]string{
				labels.Proxy: "true",
			},
		},
	}

	// Expected
	ids := []string{"testing-restart", "testing-restart-hostname"}

	// Act
	cmd := New(mock, spyOutputer{})
	err := cmd.RunE(cmd, os.Args)
	if err != nil {
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
	mock := newMockDockerClient(nil, nil, nil)

	// Act
	cmd := New(mock, spyOutputer{})
	err := cmd.RunE(cmd, os.Args)

	if err == nil {
		t.Errorf("expected the error to not be nil")
	}

	// Assert
	if len(mock.containerRestartRequests) != 0 {
		t.Errorf("expected the number of restart requests to be zero, got %d instead", len(mock.containerRestartRequests))
	}
}
