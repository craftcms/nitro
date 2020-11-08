package client

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
)

func TestInit(t *testing.T) {
	mock := newMockDockerClient(nil)
	mock.networkCreateResponse = types.NetworkCreateResponse{ID: "test-network"}

	cli := Client{
		docker: mock,
	}

	err := cli.Init(context.TODO(), "test", []string{})
	if err == nil {
		t.Errorf("expected the error to be nil, got %v", err)
	}
}
