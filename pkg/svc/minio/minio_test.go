package minio

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func TestVerifyCreated(t *testing.T) {
	type args struct {
		ctx       context.Context
		spy       *mockClient
		networkID string
		output    terminal.Outputer
	}
	tests := []struct {
		name string
		args args

		customEnvs map[string]string

		// spys
		wantSpyContainerListOptions  types.ContainerListOptions
		wantSpyImagePullImage        string
		wantSpyImagePullOptions      types.ImagePullOptions
		wantSpyContainerCreateConfig types.ContainerCreateConfig
		wantSpyContainerCreateHost   string
		wantSpyContainerStartID      string
		wantSpyContainerStartOptions types.ContainerStartOptions

		// response
		wantID       string
		wantHostname string
		wantErr      bool
	}{
		{
			name: "container is created when it does not exist",
			args: args{
				ctx: context.Background(),
				spy: &mockClient{
					// todo add image pull
					containerCreateResponse: container.ContainerCreateCreatedBody{
						ID: "someid",
					},
				},
				networkID: "some-network-id",
			},
			wantSpyContainerListOptions: types.ContainerListOptions{
				All: true,
				Filters: filters.NewArgs(
					filters.KeyValuePair{Key: "label", Value: labels.Nitro + "=true"},
					filters.KeyValuePair{Key: "label", Value: labels.Type + "=minio"},
				),
			},
			wantSpyImagePullImage: "docker.io/minio/minio:latest",
			wantSpyContainerCreateConfig: types.ContainerCreateConfig{
				Name: "minio.service.nitro",
				Config: &container.Config{
					Image: "docker.io/minio/minio:latest",
					Labels: map[string]string{
						labels.Nitro: "true",
						labels.Type:  "minio",
					},
					ExposedPorts: nat.PortSet{
						"9000/tcp": struct{}{},
					},
					Cmd: []string{"server", "/data"},
					Env: []string{"MINIO_ROOT_USER=nitro", "MINIO_ROOT_PASSWORD=nitropassword"},
				},
				HostConfig: &container.HostConfig{
					PortBindings: map[nat.Port][]nat.PortBinding{
						"9000/tcp": {
							{
								HostIP:   "127.0.0.1",
								HostPort: "9000",
							},
						},
					},
				},
				NetworkingConfig: &network.NetworkingConfig{
					EndpointsConfig: map[string]*network.EndpointSettings{
						"nitro-network": {
							NetworkID: "some-network-id",
						},
					},
				},
			},
			wantSpyContainerStartID: "someid",
			wantID:                  "someid",
			wantHostname:            "minio.service.nitro",
			wantErr:                 false,
		},
		{
			name: "custom ports are used when the environment variables are set",
			args: args{
				ctx: context.Background(),
				spy: &mockClient{
					containerCreateResponse: container.ContainerCreateCreatedBody{
						ID: "someid",
					},
				},
				networkID: "some-network-id",
			},
			customEnvs: map[string]string{
				"NITRO_MINIO_HTTP_PORT": "9001",
			},
			wantSpyContainerListOptions: types.ContainerListOptions{
				All: true,
				Filters: filters.NewArgs(
					filters.KeyValuePair{Key: "label", Value: labels.Nitro + "=true"},
					filters.KeyValuePair{Key: "label", Value: labels.Type + "=minio"},
				),
			},
			wantSpyImagePullImage: "docker.io/minio/minio:latest",
			wantSpyContainerCreateConfig: types.ContainerCreateConfig{
				Name: "minio.service.nitro",
				Config: &container.Config{
					Image: "docker.io/minio/minio:latest",
					Labels: map[string]string{
						labels.Nitro: "true",
						labels.Type:  "minio",
					},
					ExposedPorts: nat.PortSet{
						"9001/tcp": struct{}{},
					},
					Cmd: []string{"server", "/data"},
					Env: []string{"MINIO_ROOT_USER=nitro", "MINIO_ROOT_PASSWORD=nitropassword"},
				},
				HostConfig: &container.HostConfig{
					PortBindings: map[nat.Port][]nat.PortBinding{
						"9001/tcp": {
							{
								HostIP:   "127.0.0.1",
								HostPort: "9000",
							},
						},
					},
				},
				NetworkingConfig: &network.NetworkingConfig{
					EndpointsConfig: map[string]*network.EndpointSettings{
						"nitro-network": {
							NetworkID: "some-network-id",
						},
					},
				},
			},
			wantSpyContainerStartID: "someid",
			wantID:                  "someid",
			wantHostname:            "minio.service.nitro",
			wantErr:                 false,
		},
		{
			name: "containers that are already created are started",
			args: args{
				ctx: context.Background(),
				spy: &mockClient{
					containers: []types.Container{
						{
							ID:    "existing-container-id",
							State: "not-running",
						},
					},
				},
				networkID: "some-network-id",
			},
			wantSpyContainerListOptions: types.ContainerListOptions{
				All: true,
				Filters: filters.NewArgs(
					filters.KeyValuePair{Key: "label", Value: labels.Nitro + "=true"},
					filters.KeyValuePair{Key: "label", Value: labels.Type + "=minio"},
				),
			},
			wantSpyContainerStartID: "existing-container-id",
			wantID:                  "existing-container-id",
			wantHostname:            "minio.service.nitro",
			wantErr:                 false,
		},
		{
			name: "error on container list returns error",
			args: args{
				ctx: context.Background(),
				spy: &mockClient{
					containerListError: fmt.Errorf("unknown error"),
				},
			},
			wantSpyContainerListOptions: types.ContainerListOptions{
				All: true,
				Filters: filters.NewArgs(
					filters.KeyValuePair{Key: "label", Value: labels.Nitro + "=true"},
					filters.KeyValuePair{Key: "label", Value: labels.Type + "=minio"},
				),
			},
			wantID:       "",
			wantHostname: "",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		// set any custom envs
		for k, v := range tt.customEnvs {
			os.Setenv(k, v)
			defer os.Unsetenv(k)
		}

		t.Run(tt.name, func(t *testing.T) {
			id, hostname, err := VerifyCreated(tt.args.ctx, tt.args.spy, tt.args.networkID, tt.args.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyCreated() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if id != tt.wantID {
				t.Errorf("VerifyCreated() got = %v, want %v", id, tt.wantID)
			}
			if hostname != tt.wantHostname {
				t.Errorf("VerifyCreated() got1 = %v, want %v", hostname, tt.wantHostname)
			}

			// spy checks

			// check the container remove options
			if !reflect.DeepEqual(tt.wantSpyContainerListOptions, tt.args.spy.containerListOptions) {
				t.Errorf("expected the container list options to to match, got %v want %v", tt.args.spy.containerListOptions, tt.wantSpyContainerListOptions)
			}

			if tt.wantSpyImagePullImage != tt.args.spy.imagePullImage {
				t.Errorf("expected the image pull images to match, got %s want %s", tt.args.spy.imagePullImage, tt.wantSpyImagePullImage)
			}

			if !reflect.DeepEqual(tt.wantSpyContainerCreateConfig, tt.args.spy.containerCreateConfig) {
				t.Errorf("expected the container create config to to match, got %v want %v", tt.args.spy.containerCreateConfig, tt.wantSpyContainerCreateConfig)
			}

			if tt.wantSpyContainerStartID != tt.args.spy.containerStartID {
				t.Errorf("expected the container start ids to match, got %s want %s", tt.args.spy.containerStartID, tt.wantSpyContainerStartID)
			}

			if !reflect.DeepEqual(tt.wantSpyContainerStartOptions, tt.args.spy.containerStartOptions) {
				t.Errorf("expected the container start options to to match, got %v want %v", tt.args.spy.containerCreateConfig, tt.wantSpyContainerCreateConfig)
			}
		})
	}
}

func TestVerifyRemoved(t *testing.T) {
	type args struct {
		ctx    context.Context
		spy    *mockClient
		output terminal.Outputer
	}
	tests := []struct {
		name                       string
		args                       args
		wantContainerStopID        string
		wantContainerRemoveID      string
		wantContainerRemoveOptions types.ContainerRemoveOptions
		wantErr                    bool
	}{
		{
			name: "stops and removes containers when they are present and running",
			args: args{
				ctx: context.TODO(),
				spy: &mockClient{
					containers: []types.Container{
						{
							ID:    "some-random-id",
							State: "running",
						},
					},
				},
			},
			wantContainerStopID:        "some-random-id",
			wantContainerRemoveID:      "some-random-id",
			wantContainerRemoveOptions: types.ContainerRemoveOptions{RemoveVolumes: true},
			wantErr:                    false,
		},
		{
			name: "container stop returns error",
			args: args{
				ctx: context.TODO(),
				spy: &mockClient{
					containers: []types.Container{
						{
							ID:    "some-random-id",
							State: "running",
						},
					},
					containerStopError: fmt.Errorf("docker container stop error"),
				},
			},
			wantContainerStopID: "some-random-id",
			wantErr:             true,
		},
		{
			name: "container remove returns error",
			args: args{
				ctx: context.TODO(),
				spy: &mockClient{
					containers: []types.Container{
						{
							ID:    "some-random-id",
							State: "running",
						},
					},
					containerRemoveError: fmt.Errorf("docker container remove error"),
				},
			},
			wantContainerStopID:        "some-random-id",
			wantContainerRemoveID:      "some-random-id",
			wantContainerRemoveOptions: types.ContainerRemoveOptions{RemoveVolumes: true},
			wantErr:                    true,
		},
		{
			name: "non running containers do not get a stop request",
			args: args{
				ctx: context.TODO(),
				spy: &mockClient{
					containers: []types.Container{
						{
							ID:    "some-random-id",
							State: "anything",
						},
					},
				},
			},
			wantContainerStopID:        "",
			wantContainerRemoveID:      "some-random-id",
			wantContainerRemoveOptions: types.ContainerRemoveOptions{RemoveVolumes: true},
			wantErr:                    false,
		},
		{
			name: "returns no error when no containers are present",
			args: args{
				ctx: context.TODO(),
				spy: &mockClient{},
			},
			wantErr: false,
		},
		{
			name: "returns error when unable to get a list of containers",
			args: args{
				ctx: context.TODO(),
				spy: &mockClient{
					containerListError: fmt.Errorf("mock error"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// check for the error
			if err := VerifyRemoved(tt.args.ctx, tt.args.spy, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("VerifyRemoved() error = %v, wantErr %v", err, tt.wantErr)
			}

			// check the container stop id
			if tt.wantContainerStopID != "" {
				if tt.wantContainerStopID != tt.args.spy.containerStopID {
					t.Errorf("expected the container stop ids to match, got %s want %s", tt.args.spy.containerStopID, tt.wantContainerStopID)
				}
			}

			// check the container remove id
			if tt.wantContainerRemoveID != "" {
				if tt.wantContainerRemoveID != tt.args.spy.containerRemoveID {
					t.Errorf("expected the container remove ids to match, got %s want %s", tt.args.spy.containerRemoveID, tt.wantContainerRemoveID)
				}
			}

			// check the container remove options
			if !reflect.DeepEqual(tt.wantContainerRemoveOptions, tt.args.spy.containerRemoveOptions) {
				t.Errorf("expected the container remove options to to match, got %v want %v", tt.args.spy.containerRemoveOptions, tt.wantContainerRemoveOptions)
			}
		})
	}
}

type mockClient struct {
	client.CommonAPIClient

	// filters are the filters passed to list funcs
	filterArgs []filters.Args

	// mock storage
	containers           []types.Container
	containerListOptions types.ContainerListOptions
	containerListError   error

	// container create
	containerCreateConfig   types.ContainerCreateConfig
	containerCreateResponse container.ContainerCreateCreatedBody
	containerCreateError    error

	// mock start
	containerStartID      string
	containerStartOptions types.ContainerStartOptions
	containerStartError   error

	// mock stop
	containerStopID    string
	containerStopError error

	// mock remove
	containerRemoveID      string
	containerRemoveOptions types.ContainerRemoveOptions
	containerRemoveError   error

	// image pull
	imagePullReaderCloser io.ReadCloser
	imagePullImage        string
	imagePullOptions      types.ImagePullOptions
	imagePullError        error
}

func (c *mockClient) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	c.filterArgs = append(c.filterArgs, options.Filters)
	c.containerListOptions = options

	return c.containers, c.containerListError
}

func (c *mockClient) ContainerRemove(ctx context.Context, containerID string, opts types.ContainerRemoveOptions) error {
	c.containerRemoveID = containerID
	c.containerRemoveOptions = opts

	return c.containerRemoveError
}

func (c *mockClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (container.ContainerCreateCreatedBody, error) {
	c.containerCreateConfig = types.ContainerCreateConfig{
		Name:             containerName,
		Config:           config,
		HostConfig:       hostConfig,
		NetworkingConfig: networkingConfig,
	}

	return c.containerCreateResponse, c.containerCreateError
}

func (c *mockClient) ContainerStart(ctx context.Context, container string, options types.ContainerStartOptions) error {
	c.containerStartID = container
	c.containerStartOptions = options

	return c.containerStartError
}

func (c *mockClient) ContainerStop(ctx context.Context, containerID string, timeout *time.Duration) error {
	c.containerStopID = containerID

	return c.containerStopError
}

// func (c *mockClient) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
// 	// TODO(jasonmccallister) remove this hacked method
// 	summary := []types.ImageSummary{
// 		{
// 			Containers: 1,
// 		},
// 	}

// 	return summary, nil
// }

func (c *mockClient) ImagePull(ctx context.Context, image string, opts types.ImagePullOptions) (io.ReadCloser, error) {
	c.imagePullOptions = opts
	c.imagePullImage = image

	if c.imagePullReaderCloser == nil {
		c.imagePullReaderCloser = ioutil.NopCloser(ioutil.NopCloser(strings.NewReader("")))
	}

	return c.imagePullReaderCloser, c.imagePullError
}
