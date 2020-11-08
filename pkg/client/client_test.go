package client

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

// inspired by the following from the Docker docker package: https://github.com/moby/moby/blob/master/client/network_create_test.go
func newMockDockerClient(networks []types.NetworkResource, containers []types.Container, volumes []*types.Volume) *mockDockerClient {
	return &mockDockerClient{
		networkResources: networks,
		containers:       containers,
		volumes:          volumetypes.VolumesListOKBody{Volumes: volumes},
	}
}

type mockDockerClient struct {
	client.CommonAPIClient

	// container related resources for mocking calls to the client
	// the fields ending in *Response are designed to capture the
	// requests sent to the client API.
	containers              []types.Container
	containerCreateResponse container.ContainerCreateCreatedBody

	// network related resources for mocking the calls to the client
	// for network specific resources
	networkResources      []types.NetworkResource
	networkCreateResponse types.NetworkCreateResponse

	// volume related resources
	volumes              volumetypes.VolumesListOKBody
	volumeCreateResponse types.Volume

	// mockError allows us to override any func to return a method, we do not
	// set the error by default.
	mockError error
}

func (c mockDockerClient) NetworkList(ctx context.Context, options types.NetworkListOptions) ([]types.NetworkResource, error) {
	return c.networkResources, c.mockError
}

func (c mockDockerClient) NetworkCreate(ctx context.Context, name string, options types.NetworkCreate) (types.NetworkCreateResponse, error) {
	return c.networkCreateResponse, c.mockError
}

func (c mockDockerClient) VolumeList(ctx context.Context, filter filters.Args) (volumetypes.VolumesListOKBody, error) {
	return c.volumes, c.mockError
}

func (c mockDockerClient) VolumeCreate(ctx context.Context, options volumetypes.VolumesCreateBody) (types.Volume, error) {
	return c.volumeCreateResponse, c.mockError
}