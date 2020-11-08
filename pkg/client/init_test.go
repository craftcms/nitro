package client

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
)

func TestInit(t *testing.T) {
	// Arrange
	mock := newMockDockerClient(nil, nil, nil)
	mock.networkCreateResponse = types.NetworkCreateResponse{ID: "test"}
	cli := Client{docker: mock}

	// Act
	err := cli.Init(context.TODO(), "test", []string{})

	// Assert
	if err == nil {
		t.Errorf("expected the error to be nil, got %v", err)
	}
}
