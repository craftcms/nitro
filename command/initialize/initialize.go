package initialize

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command/version"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/setup"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # create a new environment with the default environment
  nitro init

  # create a new environment overriding the default name
  nitro init --environment my-new-env

  # you can override the environment by setting the variable "NITRO_DEFAULT_ENVIRONMENT"`

// NewCommand takes a docker client and returns the init command for creating a new environment
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Create new environment",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			// check if there is a config file
			_, err := config.Load(home)
			if errors.Is(err, config.ErrNoConfigFile) {
				if err := setup.FirstTime(home, output); err != nil {
					return err
				}
			}

			output.Info("Checking Nitro...")

			// create filters for the development environment
			filter := filters.NewArgs()
			filter.Add("name", "nitro-network")

			// check if the network needs to be created
			networks, err := docker.NetworkList(ctx, types.NetworkListOptions{Filters: filter})
			if err != nil {
				return fmt.Errorf("unable to list the docker networks, %w", err)
			}

			// since the filter is fuzzy, do an exact match (e.g. filtering for
			// `nitro-dev` will also return `nitro-dev-host`
			var skipNetwork bool
			var networkID string
			for _, n := range networks {
				if n.Name == "nitro-network" || strings.TrimLeft(n.Name, "/") == "nitro-network" {
					skipNetwork = true
					networkID = n.ID
				}
			}

			// create the network needs to be created
			switch skipNetwork {
			case true:
				output.Success("network ready")
			default:
				output.Pending("creating network")

				resp, err := docker.NetworkCreate(ctx, "nitro-network", types.NetworkCreate{
					Driver:     "bridge",
					Attachable: true,
					Labels: map[string]string{
						labels.Nitro:   "true",
						labels.Network: "true",
					},
				})
				if err != nil {
					return fmt.Errorf("unable to create the network, %w", err)
				}

				// set the newly created network
				networkID = resp.ID

				output.Done()
			}

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
						labels.Nitro:  "true",
						labels.Volume: "nitro",
					},
				})
				if err != nil {
					return fmt.Errorf("unable to create the volume, %w", err)
				}

				volume = &resp

				output.Done()
			}

			// build the proxy image ref
			proxyImage := fmt.Sprintf("craftcms/nitro-proxy:%s", version.Version)

			// TODO(jasonmccallister) remove this after development
			if os.Getenv("NITRO_DEVELOPMENT") != "true" {
				imageFilter := filters.NewArgs()
				imageFilter.Add("label", labels.Nitro+"=true")
				imageFilter.Add("reference", proxyImage)

				// check for the proxy image
				images, err := docker.ImageList(cmd.Context(), types.ImageListOptions{
					Filters: imageFilter,
				})
				if err != nil {
					return fmt.Errorf("unable to get a list of images, %w", err)
				}

				// if there are no local images, pull it
				if len(images) == 0 {
					output.Pending("pulling image")

					rdr, err := docker.ImagePull(ctx, proxyImage, types.ImagePullOptions{All: false})
					if err != nil {
						return fmt.Errorf("unable to pull the nitro-proxy from docker hub, %w", err)
					}

					buf := &bytes.Buffer{}
					if _, err := buf.ReadFrom(rdr); err != nil {
						return fmt.Errorf("unable to read the output from pulling the image, %w", err)
					}

					output.Done()
				}
			}

			// create a filter for the nitro proxy
			proxyFilter := filters.NewArgs()
			proxyFilter.Add("label", labels.Nitro+"=true")
			proxyFilter.Add("label", labels.Proxy+"=true")

			// check if there is an existing container for the nitro-proxy
			var containerID string
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: proxyFilter, All: true})
			if err != nil {
				return fmt.Errorf("unable to list the containers\n%w", err)
			}

			var proxyRunning bool
			for _, c := range containers {
				for _, n := range c.Names {
					if n == "nitro-proxy" || n == "/nitro-proxy" {
						output.Success("proxy ready")

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
				output.Pending("creating proxy")

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
					httpPort, err = nat.NewPort("tcp", os.Getenv("NITRO_HTTP_PORT"))
					if err != nil {
						return fmt.Errorf("unable to set the HTTP port, %w", err)
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
					httpsPort, _ = nat.NewPort("tcp", os.Getenv("NITRO_HTTPS_PORT"))
					if err != nil {
						return fmt.Errorf("unable to set the HTTPS port, %w", err)
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
					apiPort, err = nat.NewPort("tcp", os.Getenv("NITRO_API_PORT"))
					if err != nil {
						return fmt.Errorf("unable to set the API port, %w", err)
					}
				}

				// create a container
				resp, err := docker.ContainerCreate(ctx,
					&container.Config{
						Image: proxyImage,
						ExposedPorts: nat.PortSet{
							httpPort:  struct{}{},
							httpsPort: struct{}{},
							apiPort:   struct{}{},
						},
						Labels: map[string]string{
							labels.Nitro:        "true",
							labels.Type:         "proxy",
							labels.Proxy:        "true",
							labels.ProxyVersion: version.Version,
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
							"nitro-network": {
								NetworkID: networkID,
							},
						},
					},
					nil,
					"nitro-proxy",
				)
				if err != nil {
					return fmt.Errorf("unable to create the container from image %s\n%w", proxyImage, err)
				}

				containerID = resp.ID

				output.Done()
			}

			// start the container for the proxy if its not running
			if !proxyRunning {
				if err := docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
					return fmt.Errorf("unable to start the nitro container, %w", err)
				}
			}

			// convert the apply flag to a boolean
			skipApply, err := strconv.ParseBool(cmd.Flag("skip-apply").Value.String())
			if err != nil {
				// don't do anything
			}

			// check if we need to run the
			if skipApply != true && cmd.Parent() != nil {
				// TODO(jasonmccallister) make this better :)
				for _, c := range cmd.Parent().Commands() {
					// set the apply command
					if c.Use == "apply" {
						if err := c.RunE(c, args); err != nil {
							return err
						}
					}

					// set the trust command
					if c.Use == "trust" {
						if err := c.RunE(c, args); err != nil {
							return err
						}
					}
				}
			}

			output.Info("Nitro is ready! 🚀")

			return nil
		},
	}

	// set flags for the command
	cmd.Flags().Bool("skip-apply", false, "skip applying changes")

	return cmd
}
