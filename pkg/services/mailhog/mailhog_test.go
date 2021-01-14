package mailhog

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

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
	containers         []types.Container
	containerListError error

	// mock stop
	containerStopID    string
	containerStopError error

	// mock remove
	containerRemoveID      string
	containerRemoveOptions types.ContainerRemoveOptions
	containerRemoveError   error
}

func (c *mockClient) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	c.filterArgs = append(c.filterArgs, options.Filters)

	return c.containers, c.containerListError
}

func (c *mockClient) ContainerRemove(ctx context.Context, containerID string, opts types.ContainerRemoveOptions) error {
	c.containerRemoveID = containerID
	c.containerRemoveOptions = opts

	return c.containerRemoveError
}

// func (c *mockClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (container.ContainerCreateCreatedBody, error) {
// 	// save the request on the struct field
// 	// TODO(jasonmccallister) this is wrong, need to look at the code to determine the correct
// 	// types are set and returned
// 	c.containerCreateRequests = append(c.containerCreateRequests, types.ContainerCreateConfig{
// 		Name:             containerName,
// 		Config:           config,
// 		HostConfig:       hostConfig,
// 		NetworkingConfig: networkingConfig,
// 	})

// 	return c.containerCreateResponse, c.mockError
// }

// func (c *mockClient) ContainerStart(ctx context.Context, container string, options types.ContainerStartOptions) error {
// 	c.containerID = container
// 	c.containerStartRequests = append(c.containerStartRequests, options)

// 	return c.mockError
// }

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
