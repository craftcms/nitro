package initialize

import (
	"os"
	"reflect"
	"testing"

	"github.com/craftcms/nitro/labels"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/go-connections/nat"
)

func TestInitFromFreshCreatesNewResources(t *testing.T) {
	// Arrange
	environmentName := "testing-init"
	mock := newMockDockerClient(nil, nil, nil)
	mock.networkCreateResponse = types.NetworkCreateResponse{
		ID: "testing-init",
	}
	mock.containerCreateResponse = container.ContainerCreateCreatedBody{
		ID: "testingid",
	}

	// Expected
	// set the network create request
	networkReq := types.NetworkCreateRequest{
		NetworkCreate: types.NetworkCreate{
			Driver:     "bridge",
			Attachable: true,
			Labels: map[string]string{
				labels.Environment: "testing-init",
				labels.Network:     "testing-init",
			},
		},
		Name: "testing-init",
	}
	// set the volume create request
	volumeReq := volumetypes.VolumesCreateBody{
		Driver: "local",
		Name:   "testing-init",
		Labels: map[string]string{
			labels.Environment: "testing-init",
			labels.Volume:      "testing-init",
		},
	}
	// set the container create request
	containerCreateReq := types.ContainerCreateConfig{
		Config: &container.Config{
			Image: "craftcms/nitro-proxy:develop",
			ExposedPorts: nat.PortSet{
				"80/tcp":   struct{}{},
				"443/tcp":  struct{}{},
				"5000/tcp": struct{}{},
			},
			Labels: map[string]string{
				labels.Type:         "proxy",
				labels.Environment:  "testing-init",
				labels.Proxy:        "testing-init",
				labels.ProxyVersion: "develop",
			},
		},
		HostConfig: &container.HostConfig{
			NetworkMode: "default",
			Mounts: []mount.Mount{
				{
					Type: mount.TypeVolume,
					// TODO(jasonmccallister) fix the mock to return, or filter, volumes
					Source: "",
					Target: "/data",
				},
			},
			PortBindings: map[nat.Port][]nat.PortBinding{
				"80/tcp": {
					{
						HostIP:   "127.0.0.1",
						HostPort: "80",
					},
				},
				"443/tcp": {
					{
						HostIP:   "127.0.0.1",
						HostPort: "443",
					},
				},
				"5000/tcp": {
					{
						HostIP:   "127.0.0.1",
						HostPort: "5000",
					},
				},
			},
		},
		NetworkingConfig: &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"testing-init": {
					NetworkID: "testing-init",
				},
			},
		},
		Name: "testing-init",
	}
	// set the container start request
	containerStartRequest := types.ContainerStartOptions{}

	// Act
	cmd := NewCommand(mock, spyOutputer{})
	cmd.Flags().StringP("environment", "e", environmentName, "test flag")
	err := cmd.RunE(cmd, os.Args)

	// Assert
	if err != nil {
		t.Errorf("expected the error to be nil, got %v", err)
	}

	// make sure the network create matches the expected
	if !reflect.DeepEqual(mock.networkCreateRequests[0], networkReq) {
		t.Errorf(
			"expected network create request to match\ngot:\n%v\nwant:\n%v",
			mock.networkCreateRequests[0],
			networkReq,
		)
	}

	// make sure the volume create matches the expected
	if !reflect.DeepEqual(mock.volumeCreateRequest, volumeReq) {
		t.Errorf(
			"expected volume create request to match\ngot:\n%v\nwant:\n%v",
			mock.volumeCreateRequest,
			volumeReq,
		)
	}

	// make sure the container create matches the expected
	if !reflect.DeepEqual(mock.containerCreateRequests[0], containerCreateReq) {
		if !reflect.DeepEqual(mock.containerCreateRequests[0].Config, containerCreateReq.Config) {
			t.Errorf(
				"expected container create request config to match\ngot:\n%v\n\nwant:\n%v",
				mock.containerCreateRequests[0].Config,
				containerCreateReq.Config,
			)
		}

		if !reflect.DeepEqual(mock.containerCreateRequests[0].HostConfig, containerCreateReq.HostConfig) {
			t.Errorf(
				"expected container create request host config to match\ngot:\n%v\n\nwant:\n%v",
				mock.containerCreateRequests[0].HostConfig,
				containerCreateReq.HostConfig,
			)
		}

		if !reflect.DeepEqual(mock.containerCreateRequests[0].NetworkingConfig, containerCreateReq.NetworkingConfig) {
			t.Errorf(
				"expected container create request networking to match\ngot:\n%v\n\nwant:\n%v",
				mock.containerCreateRequests[0].NetworkingConfig,
				containerCreateReq.NetworkingConfig,
			)
		}
	}

	// make sure the container start matches the expected
	if !reflect.DeepEqual(mock.containerStartRequests[0], containerStartRequest) {
		t.Errorf(
			"expected container start request to match\ngot:\n%v\nwant:\n%v",
			mock.containerStartRequests[0],
			containerStartRequest,
		)
	}

	// make sure the container ID to start matches the expected
	if mock.containerID != "testingid" {
		t.Errorf(
			"expected container IDs to start to match\ngot:\n%v\nwant:\n%v",
			mock.containerID,
			"testingid",
		)
	}
}
