package proxycontainer

import (
	"bytes"
	"context"
	"fmt"
	"os"

	volumetypes "github.com/docker/docker/api/types/volume"

	"github.com/craftcms/nitro/command/version"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

var (
	// ProxyImage is the docker hub image with the current CLI version
	ProxyImage = fmt.Sprintf("craftcms/nitro-proxy:%s", version.Version)

	// ProxyName is the name of the proxy container (e.g. nitro-proxy)
	ProxyName = "nitro-proxy"

	// ErrNoProxyContainer is returned when the proxy container is not found
	ErrNoProxyContainer = fmt.Errorf("unable to locate the proxy container")
)

// Create is used to create a new proxy container for the nitro development environment.
func Create(ctx context.Context, docker client.CommonAPIClient, output terminal.Outputer, networkID string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	filter := filters.NewArgs()
	filter.Add("label", containerlabels.Nitro+"=true")
	filter.Add("reference", ProxyImage)

	// check for the proxy image
	images, err := docker.ImageList(ctx, types.ImageListOptions{Filters: filter})
	if err != nil {
		return fmt.Errorf("unable to get a list of images, %w", err)
	}

	// if there are no local images, pull it
	if len(images) == 0 && os.Getenv("NITRO_DEVELOPMENT") != "true" {
		output.Pending("pulling image")

		rdr, err := docker.ImagePull(ctx, ProxyImage, types.ImagePullOptions{All: false})
		if err != nil {
			return fmt.Errorf("unable to pull the nitro-proxy from docker hub, %w", err)
		}

		buf := &bytes.Buffer{}
		if _, err := buf.ReadFrom(rdr); err != nil {
			return fmt.Errorf("unable to read the output from pulling the image, %w", err)
		}

		output.Done()
	}

	filter.Del("reference", ProxyImage)
	// check if the volume needs to be created
	volumes, err := docker.VolumeList(ctx, filter)
	if err != nil {
		return fmt.Errorf("unable to list volumes, %w", err)
	}

	// since the filter is fuzzy, do an exact match (e.g. filtering for
	// `nitro-dev` will also return `nitro-dev-host`
	var skipVolume bool
	var volume *types.Volume
	for _, v := range volumes.Volumes {
		if v.Name == "nitro" {
			skipVolume = true
			volume = v
		}
	}

	// check if the volume needs to be created
	switch skipVolume {
	case true:
		output.Success("volume ready")
	default:
		output.Pending("creating volume")

		// create a volume with the same name of the machine
		resp, err := docker.VolumeCreate(ctx, volumetypes.VolumeCreateBody{
			Driver: "local",
			Name:   "nitro",
			Labels: map[string]string{
				containerlabels.Nitro:  "true",
				containerlabels.Volume: "nitro",
			},
		})
		if err != nil {
			return fmt.Errorf("unable to create the volume, %w", err)
		}

		volume = &resp

		output.Done()
	}

	// remove the reference filter
	filter.Del("reference", ProxyImage)

	// create a filter for the nitro proxy
	filter.Add("label", containerlabels.Proxy+"=true")

	// check if there is an existing container for the nitro-proxy
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter, All: true})
	if err != nil {
		return fmt.Errorf("unable to list the containers\n%w", err)
	}

	// check the containers and verify its running
	for _, c := range containers {
		for _, n := range c.Names {
			if n == "nitro-proxy" || n == "/nitro-proxy" {
				// check if it is running
				if c.State != "running" {
					if err := docker.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
						return fmt.Errorf("unable to start the nitro container, %w", err)
					}
				}

				output.Success("proxy ready")

				return nil
			}
		}
	}

	// if we do not have a proxy, it needs to be create
	output.Pending("creating proxy")

	// check for a custom HTTP port
	httpPort := "80"
	if _, defined := os.LookupEnv("NITRO_HTTP_PORT"); defined {
		httpPort = os.Getenv("NITRO_HTTP_PORT")
	}

	// check for a custom HTTPS port
	httpsPort := "443"
	if _, defined := os.LookupEnv("NITRO_HTTPS_PORT"); defined {
		httpsPort = os.Getenv("NITRO_HTTPS_PORT")
	}

	// check for a custom API port
	apiPort := "5000"
	if _, defined := os.LookupEnv("NITRO_API_PORT"); defined {
		apiPort = os.Getenv("NITRO_API_PORT")
	}

	httpPortNat, err := nat.NewPort("tcp", "80")
	if err != nil {
		return fmt.Errorf("unable to set the HTTP port, %w", err)
	}

	httpsPortNat, err := nat.NewPort("tcp", "443")
	if err != nil {
		return fmt.Errorf("unable to set the HTTPS port, %w", err)
	}

	apiPortNat, err := nat.NewPort("tcp", "5000")
	if err != nil {
		return fmt.Errorf("unable to set the API port, %w", err)
	}

	// create a container
	resp, err := docker.ContainerCreate(ctx,
		&container.Config{
			Image: ProxyImage,
			ExposedPorts: nat.PortSet{
				httpPortNat:  struct{}{},
				httpsPortNat: struct{}{},
				apiPortNat:   struct{}{},
			},
			Labels: map[string]string{
				containerlabels.Nitro:        "true",
				containerlabels.Type:         "proxy",
				containerlabels.Proxy:        "true",
				containerlabels.ProxyVersion: version.Version,
			},
			Env: []string{"PGPASSWORD=nitro", "PGUSER=nitro", "NITRO_VERSION=" + version.Version},
		},
		&container.HostConfig{
			NetworkMode: "default",
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: volume.Name,
					Target: "/data",
				},
			},
			PortBindings: map[nat.Port][]nat.PortBinding{
				httpPortNat: {
					{
						HostIP:   "127.0.0.1",
						HostPort: httpPort,
					},
				},
				httpsPortNat: {
					{
						HostIP:   "127.0.0.1",
						HostPort: httpsPort,
					},
				},
				apiPortNat: {
					{
						HostIP:   "127.0.0.1",
						HostPort: apiPort,
					},
				},
			},
		},
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"nitro-network": {
					NetworkID: networkID,
				},
			},
		},
		nil,
		ProxyName,
	)
	if err != nil {
		return fmt.Errorf("unable to create proxy container: %s\n%w", ProxyImage, err)
	}

	if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("unable to start the nitro container, %w", err)
	}

	output.Done()

	return nil
}

// FindAndStart will look for the proxy container and verify the container is started. It will return the
// ErrNoProxyContainer error if it is unable to locate the proxy container. It is NOT responsible for
// creating the proxy container as that is handled in the initialize package.
func FindAndStart(ctx context.Context, docker client.ContainerAPIClient) (types.Container, error) {
	// create the filters for the proxy
	f := filters.NewArgs()
	f.Add("label", containerlabels.Type+"=proxy")

	// check if there is an existing container for the nitro-proxy
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: f, All: true})
	if err != nil {
		return types.Container{}, fmt.Errorf("unable to list the containers: %w", err)
	}

	for _, c := range containers {
		for _, n := range c.Names {
			if n == "nitro-proxy" || n == "/nitro-proxy" {
				// check if it is running
				if c.State != "running" {
					if err := docker.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
						return types.Container{}, fmt.Errorf("unable to start the proxy container: %w", err)
					}
				}

				// return the container
				return c, nil
			}
		}
	}

	return types.Container{}, ErrNoProxyContainer
}
