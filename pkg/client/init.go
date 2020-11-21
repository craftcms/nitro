package client

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/go-connections/nat"
)

// Init is responsible for creating the networks, volumes, and proxy containers for the
// nitro docker setup. It will check for the existing networks, volumes, and containers
// instead of overwriting the containers. Init should only be used once to setup the
// development environment, which is why we safeguard the resources.
func (cli *Client) Init(ctx context.Context, name string, args []string) error {
	cli.Info(fmt.Sprintf("Checking %s...", name))

	// create filters for the development environment
	filter := filters.NewArgs()
	filter.Add("name", name)

	// check if the network needs to be created
	networks, err := cli.docker.NetworkList(ctx, types.NetworkListOptions{Filters: filter})
	if err != nil {
		return fmt.Errorf("unable to list the docker networks, %w", err)
	}

	// since the filter is fuzzy, do an exact match (e.g. filtering for
	// `nitro-dev` will also return `nitro-dev-host`
	var skipNetwork bool
	var networkID string
	for _, n := range networks {
		if n.Name == name {
			skipNetwork = true
			networkID = n.ID
		}
	}

	// create the network needs to be created
	switch skipNetwork {
	case true:
		cli.InfoSuccess("network ready")
	default:
		cli.InfoPending("creating network")

		resp, err := cli.docker.NetworkCreate(ctx, name, types.NetworkCreate{
			Driver:     "bridge",
			Attachable: true,
			Labels: map[string]string{
				EnvironmentLabel: name,
				NetworkLabel:     name,
			},
		})
		if err != nil {
			return fmt.Errorf("unable to create the network, %w", err)
		}

		// set the newly created network
		networkID = resp.ID

		cli.InfoDone()
	}

	// check if the volume needs to be created
	volumes, err := cli.docker.VolumeList(ctx, filter)
	if err != nil {
		return fmt.Errorf("unable to list volumes, %w", err)
	}

	// since the filter is fuzzy, do an exact match (e.g. filtering for
	// `nitro-dev` will also return `nitro-dev-host`
	var skipVolume bool
	var volume *types.Volume
	for _, v := range volumes.Volumes {
		if v.Name == name {
			skipVolume = true
			volume = v
		}
	}

	// check if the volume needs to be created
	switch skipVolume {
	case true:
		cli.InfoSuccess("volume ready")
	default:
		cli.InfoPending("creating volume")

		// create a volume with the same name of the machine
		resp, err := cli.docker.VolumeCreate(ctx, volumetypes.VolumesCreateBody{
			Driver: "local",
			Name:   name,
			Labels: map[string]string{
				EnvironmentLabel:            name,
				"com.craftcms.nitro.volume": name,
			},
		})
		if err != nil {
			return fmt.Errorf("unable to create the volume, %w", err)
		}

		volume = &resp

		cli.InfoDone()
	}

	// pull the latest image from docker hub for the nitro-proxy
	// TODO(jasonmccallister) replace with the nitr o CLI version and non-local image (e.g. craftcms/nitro-proxy:version)
	// _, err = cli.docker.ImagePull(ctx, "nitro-proxy:develop", types.ImagePullOptions{})
	// if err != nil {
	// 	return fmt.Errorf("unable to pull the nitro-proxy from docker hub, %w", err)
	// }

	// create a filter for the nitro proxy
	pf := filters.NewArgs()
	pf.Add("label", "com.craftcms.nitro.proxy="+name)

	// check if there is an existing container for the nitro-proxy
	var containerID string
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{Filters: pf, All: true})
	if err != nil {
		return fmt.Errorf("unable to list the containers\n%w", err)
	}

	var proxyRunning bool
	for _, c := range containers {
		for _, n := range c.Names {
			if n == name || n == "/"+name {
				cli.InfoSuccess("proxy ready")

				containerID = c.ID

				// check if it is running
				if c.State == "running" {
					proxyRunning = true
				}
			}
		}
	}

	// if we do not have a container id, it needs to be create
	if containerID == "" {
		cli.InfoPending("creating proxy")

		// set ports
		var httpPort, httpsPort, apiPort nat.Port

		// check for a custom HTTP port
		switch os.Getenv("NITRO_HTTP_PORT") {
		case "":
			httpPort, err = nat.NewPort("tcp", "80")
			if err != nil {
				return fmt.Errorf("unable to set the HTTP port, %w", err)
			}
		default:
			if os.Getenv("NITRO_HTTP_PORT") != "" {
				httpPort, err = nat.NewPort("tcp", os.Getenv("NITRO_HTTP_PORT"))
				if err != nil {
					return fmt.Errorf("unable to set the HTTP port, %w", err)
				}
			}
		}

		// check for a custom HTTPS port
		switch os.Getenv("NITRO_HTTPS_PORT") {
		case "":
			httpsPort, err = nat.NewPort("tcp", "443")
			if err != nil {
				return fmt.Errorf("unable to set the HTTPS port, %w", err)
			}
		default:
			if os.Getenv("NITRO_HTTPS_PORT") != "" {
				httpsPort, _ = nat.NewPort("tcp", os.Getenv("NITRO_HTTPS_PORT"))
				if err != nil {
					return fmt.Errorf("unable to set the HTTPS port, %w", err)
				}
			}
		}

		// check for a custom API port
		switch os.Getenv("NITRO_API_PORT") {
		case "":
			apiPort, err = nat.NewPort("tcp", "5000")
			if err != nil {
				return fmt.Errorf("unable to set the API port, %w", err)
			}
		default:
			if os.Getenv("NITRO_API_PORT") != "" {
				httpPort, _ = nat.NewPort("tcp", os.Getenv("NITRO_API_PORT"))
				if err != nil {
					return fmt.Errorf("unable to set the API port, %w", err)
				}
			}
		}

		// create a container
		resp, err := cli.docker.ContainerCreate(ctx,
			&container.Config{
				Image: "nitro-proxy:develop",
				ExposedPorts: nat.PortSet{
					httpPort:  struct{}{},
					httpsPort: struct{}{},
					apiPort:   struct{}{},
				},
				Labels: map[string]string{
					"com.craftcms.nitro.type":  "proxy",
					EnvironmentLabel:           name,
					"com.craftcms.nitro.proxy": name,
				},
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
					httpPort: {
						{
							HostIP:   "127.0.0.1",
							HostPort: "80",
						},
					},
					httpsPort: {
						{
							HostIP:   "127.0.0.1",
							HostPort: "443",
						},
					},
					apiPort: {
						{
							HostIP:   "127.0.0.1",
							HostPort: "5000",
						},
					},
				},
			},
			&network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					name: {
						NetworkID: networkID,
					},
				},
			},
			name,
		)
		if err != nil {
			return fmt.Errorf("unable to create the container\n%w", err)
		}

		containerID = resp.ID

		cli.InfoDone()
	}

	// start the container for the proxy if its not running
	if !proxyRunning {
		if err := cli.docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
			return fmt.Errorf("unable to start the nitro container, %w", err)
		}
	}

	cli.Info(name, "is ready! 🚀")

	return nil
}
