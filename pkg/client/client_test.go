package client

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// inspired by the following from the Docker docker package: https://github.com/moby/moby/blob/master/client/network_create_test.go

func newMockDockerClient(networks []types.NetworkResource) *mockDockerClient {
	return &mockDockerClient{
		networkResources: networks,
	}
}

type mockDockerClient struct {
	client.CommonAPIClient

	networkResources      []types.NetworkResource
	networkCreateResponse types.NetworkCreateResponse
	mockError             error
}

func (c mockDockerClient) NetworkList(ctx context.Context, options types.NetworkListOptions) ([]types.NetworkResource, error) {
	return c.networkResources, c.mockError
}

func (c mockDockerClient) NetworkCreate(ctx context.Context, name string, options types.NetworkCreate) (types.NetworkCreateResponse, error) {
	return c.networkCreateResponse, c.mockError
}
