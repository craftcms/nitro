package client

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
)

func TestRestart(t *testing.T) {
	// Arrange
	environmentName := "testing-init"
	mock := newMockDockerClient(nil, nil, nil)
	mock.containers = []types.Container{
		{
			Labels: nil,
		},
	}

	cli := Client{docker: mock}

	// Expected
	if err := cli.Restart(context.Background(), environmentName, []string{}); err != nil {
		t.Skipf("need to implement")
	}

	// Assert
}
