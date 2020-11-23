package initialize

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

// inspired by the following from the Docker docker package: https://github.com/moby/moby/blob/master/client/network_create_test.go
func newMockDockerClient(networks []types.NetworkResource, containers []types.Container, volumes []*types.Volume) *mockDockerClient {
	return &mockDockerClient{
		networks:   networks,
		containers: containers,
		volumes:    volumetypes.VolumesListOKBody{Volumes: volumes},
	}
}

type mockDockerClient struct {
	client.CommonAPIClient

	// filters are the filters passed to list funcs
	filterArgs []filters.Args

	// container related resources for mocking calls to the client
	// the fields ending in *Response are designed to capture the
	// requests sent to the client API.
	containerID              string
	containers               []types.Container
	containerCreateRequests  []types.ContainerCreateConfig
	containerCreateResponse  container.ContainerCreateCreatedBody
	containerStartRequests   []types.ContainerStartOptions
	containerRestartRequests []string

	// network related resources for mocking the calls to the client
	// for network specific resources
	networks              []types.NetworkResource
	networkCreateRequests []types.NetworkCreateRequest
	networkCreateResponse types.NetworkCreateResponse

	// volume related resources
	volumes              volumetypes.VolumesListOKBody
	volumeCreateRequest  volumetypes.VolumesCreateBody
	volumeCreateResponse types.Volume

	// mockError allows us to override any func to return a method, we do not
	// set the error by default.
	mockError error
}

func (c *mockDockerClient) NetworkList(ctx context.Context, options types.NetworkListOptions) ([]types.NetworkResource, error) {
	c.filterArgs = append(c.filterArgs, options.Filters)

	return c.networks, c.mockError
}

func (c *mockDockerClient) NetworkCreate(ctx context.Context, name string, options types.NetworkCreate) (types.NetworkCreateResponse, error) {
	// save the request on the struct field
	c.networkCreateRequests = append(c.networkCreateRequests, types.NetworkCreateRequest{
		NetworkCreate: options,
		Name:          name,
	})

	return c.networkCreateResponse, c.mockError
}

func (c *mockDockerClient) VolumeList(ctx context.Context, filter filters.Args) (volumetypes.VolumesListOKBody, error) {
	c.filterArgs = append(c.filterArgs, filter)

	return c.volumes, c.mockError
}

func (c *mockDockerClient) VolumeCreate(ctx context.Context, options volumetypes.VolumesCreateBody) (types.Volume, error) {
	c.volumeCreateRequest = options

	return c.volumeCreateResponse, c.mockError
}

func (c *mockDockerClient) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	c.filterArgs = append(c.filterArgs, options.Filters)

	return c.containers, c.mockError
}

func (c *mockDockerClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error) {
	// save the request on the struct field
	// TODO(jasonmccallister) this is wrong, need to look at the code to determine the correct
	// types are set and returned
	c.containerCreateRequests = append(c.containerCreateRequests, types.ContainerCreateConfig{
		Name:             containerName,
		Config:           config,
		HostConfig:       hostConfig,
		NetworkingConfig: networkingConfig,
	})

	return c.containerCreateResponse, c.mockError
}

func (c *mockDockerClient) ContainerStart(ctx context.Context, container string, options types.ContainerStartOptions) error {
	c.containerID = container
	c.containerStartRequests = append(c.containerStartRequests, options)

	return c.mockError
}

func (c *mockDockerClient) ContainerRestart(ctx context.Context, container string, timeout *time.Duration) error {
	c.containerRestartRequests = append(c.containerRestartRequests, container)
	return c.mockError
}

func (c *mockDockerClient) ContainerStop(ctx context.Context, containerID string, timeout *time.Duration) error {
	c.containerID = containerID

	return c.mockError
}
