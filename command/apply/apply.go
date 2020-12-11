package apply

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
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

	"github.com/craftcms/nitro/command/apply/internal/match"
	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/labels"
	"github.com/craftcms/nitro/pkg/sudo"
	"github.com/craftcms/nitro/protob"
	"github.com/craftcms/nitro/terminal"
)

var (
	// ErrNoNetwork is used when we cannot find the network
	ErrNoNetwork = fmt.Errorf("Unable to find the network")

	// NginxImage is the image used for sites, with the PHP version
	NginxImage = "docker.io/craftcms/nginx:%s-dev"

	// DatabaseImage is used for determining the engine and version
	DatabaseImage = "docker.io/library/%s:%s"
)

const exampleText = `  # apply changes from a config
  nitro apply`

// New takes a docker client and the terminal output to run the apply actions
func New(home string, docker client.CommonAPIClient, nitrod protob.NitroClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apply",
		Short:   "Apply changes to an environment",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			env := cmd.Flag("environment").Value.String()
			ctx := cmd.Context()
			if ctx == nil {
				// when we call commands from other commands (e.g. init)
				// the context could be nil, so we set it to the parent
				// context just in case.
				ctx = cmd.Parent().Context()
			}

			cfg, err := config.Load(home, env)
			if err != nil {
				return err
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", labels.Environment+"="+env)

			output.Info("Checking Network...")

			// find networks
			networks, err := docker.NetworkList(ctx, types.NetworkListOptions{Filters: filter})
			if err != nil {
				return fmt.Errorf("unable to list docker networks\n%w", err)
			}

			// get the network for the environment
			var networkID string
			for _, n := range networks {
				if n.Name == env {
					networkID = n.ID
					output.Success("network ready")
				}
			}

			// if the network is not found
			if networkID == "" {
				return ErrNoNetwork
			}

			output.Info("Checking Proxy...")

			proxyFilter := filters.NewArgs()
			proxyFilter.Add("label", labels.Proxy+"="+env)

			// check if there is an existing container for the nitro-proxy
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: proxyFilter, All: true})
			if err != nil {
				return fmt.Errorf("unable to list the containers\n%w", err)
			}

			// get the container id and determine if the container needs to start
			var proxyContainerID string
			var proxyRunning bool
			for _, c := range containers {
				for _, n := range c.Names {
					if n == env || n == "/"+env {
						proxyContainerID = c.ID

						// check if it is running
						if c.State == "running" {
							output.Success("proxy ready")
							proxyRunning = true
						}
					}
				}
			}

			// if its not running start it
			if !proxyRunning {
				output.Pending("starting proxy")

				if err := docker.ContainerStart(ctx, proxyContainerID, types.ContainerStartOptions{}); err != nil {
					return fmt.Errorf("unable to start the nitro container, %w", err)
				}

				output.Done()
			}

			// check the databases
			output.Info("Checking Databases...")
			for _, db := range cfg.Databases {
				// add filters to check for the container
				filter.Add("label", labels.DatabaseEngine+"="+db.Engine)
				filter.Add("label", labels.DatabaseVersion+"="+db.Version)
				filter.Add("label", labels.Type+"=database")

				// get the containers for the database
				containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
				if err != nil {
					return fmt.Errorf("error getting a list of containers")
				}

				// set the hostname for the database container
				hostname, err := db.GetHostname()
				if err != nil {
					return err
				}

				// if there is not a container for the database, create a volume, container, and start the container
				var containerID string
				var startContainer bool
				switch len(containers) {
				// the database container exists
				case 1:
					// set the container id
					containerID = containers[0].ID

					// check if the container is running
					if containers[0].State != "running" {
						startContainer = true
						output.Pending("starting", hostname)
					} else {
						output.Success(hostname, "ready")
					}
				default:
					// database container does not exist, so create the volume and start it
					output.Pending("creating volume", hostname)

					// create the database labels
					lbls := map[string]string{
						labels.Environment:     env,
						labels.DatabaseEngine:  db.Engine,
						labels.DatabaseVersion: db.Version,
						labels.Type:            "database",
					}

					// if the database is mysql or mariadb, mark them as
					// mysql compatible (used for importing backups)
					if db.Engine == "mariadb" || db.Engine == "mysql" {
						lbls[labels.DatabaseCompatability] = "mysql"
					}

					// if the database is postgres, mark it as compatible
					// with postgres. This is not needed but a place holder
					// if cockroachdb is ever supported by craft.
					if db.Engine == "postgres" {
						lbls[labels.DatabaseCompatability] = "postgres"
					}

					// create the volume
					volResp, err := docker.VolumeCreate(ctx, volumetypes.VolumesCreateBody{
						Driver: "local",
						Name:   hostname,
						Labels: lbls,
					})
					if err != nil {
						return fmt.Errorf("unable to create the volume, %w", err)
					}

					output.Done()

					// determine the image name
					image := fmt.Sprintf(DatabaseImage, db.Engine, db.Version)

					// set mounts and environment based on the database type
					target := "/var/lib/mysql"
					var envs []string
					if strings.Contains(image, "postgres") {
						target = "/var/lib/postgresql/data"
						envs = []string{"POSTGRES_USER=nitro", "POSTGRES_DB=nitro", "POSTGRES_PASSWORD=nitro"}
					} else {
						envs = []string{"MYSQL_ROOT_PASSWORD=nitro", "MYSQL_DATABASE=nitro", "MYSQL_USER=nitro", "MYSQL_PASSWORD=nitro"}
					}

					// check for if we should skip pulling an image
					if cmd.Flag("skip-pull").Value.String() == "false" {
						output.Pending("pulling", image)

						// pull the image
						rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
						if err != nil {
							return fmt.Errorf("unable to pull image %s, %w", image, err)
						}

						// read the output to pull the image
						buf := &bytes.Buffer{}
						if _, err := buf.ReadFrom(rdr); err != nil {
							return fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
						}

						output.Done()
					}

					// set the port for the database
					port, err := nat.NewPort("tcp", db.Port)
					if err != nil {
						return fmt.Errorf("unable to create the port, %w", err)
					}

					output.Pending("creating", hostname)

					// tell docker to create the database container
					resp, err := docker.ContainerCreate(
						ctx,
						&container.Config{
							Image:  image,
							Labels: lbls,
							ExposedPorts: nat.PortSet{
								port: struct{}{},
							},
							Env: envs,
						},
						&container.HostConfig{
							Mounts: []mount.Mount{
								{
									Type:   mount.TypeVolume,
									Source: volResp.Name,
									Target: target,
								},
							},
							PortBindings: map[nat.Port][]nat.PortBinding{
								port: {
									{
										HostIP:   "127.0.0.1",
										HostPort: db.Port,
									},
								},
							},
						},
						&network.NetworkingConfig{
							EndpointsConfig: map[string]*network.EndpointSettings{
								env: {
									NetworkID: networkID,
								},
							},
						},
						hostname,
					)
					if err != nil {
						return fmt.Errorf("unable to create the container, %w", err)
					}

					// set the database container id to start
					containerID = resp.ID
					startContainer = true
				}

				// start the database container if needed
				if startContainer {
					if err := docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
						return fmt.Errorf("unable to start the container, %w", err)
					}

					output.Done()
				}

				// remove the database filters
				filter.Del("label", labels.DatabaseEngine+"="+db.Engine)
				filter.Del("label", labels.DatabaseVersion+"="+db.Version)
				filter.Del("label", labels.Type+"=database")
			}

			// get all of the sites, their local path, the php version, and the type of project (nginx or PHP-FPM)
			output.Info("Checking Sites...")

			// get the envs for the sites
			envs := cfg.AsEnvs()

			for _, site := range cfg.Sites {
				// add the site filter
				filter.Add("label", labels.Host+"="+site.Hostname)

				// look for a container for the site
				containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
				if err != nil {
					return fmt.Errorf("error getting a list of containers")
				}

				switch len(containers) {
				case 1:
					// there is a running container
					c := containers[0]
					image := fmt.Sprintf(NginxImage, site.PHP)

					// TODO(jasonmccallister) get the containers environment variables

					// make sure the images and mounts match, if they don't stop, remove, and recreate the container
					if match.Site(home, site, cfg.PHP, c) == false {
						output.Pending(site.Hostname, "out of sync")

						path, err := site.GetAbsPath(home)
						if err != nil {
							return err
						}

						// stop container
						if err := docker.ContainerStop(ctx, c.ID, nil); err != nil {
							return err
						}

						// remove container
						if err := docker.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{}); err != nil {
							return err
						}

						output.Done()

						// pull the image
						if cmd.Flag("skip-pull").Value.String() == "false" {
							output.Pending("pulling", image)

							// pull the image
							rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
							if err != nil {
								return fmt.Errorf("unable to pull image, %w", err)
							}

							// read to pull the image
							buf := &bytes.Buffer{}
							if _, err := buf.ReadFrom(rdr); err != nil {
								return fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
							}

							output.Done()
						}

						// add the path mount
						mounts := []mount.Mount{}
						mounts = append(mounts, mount.Mount{
							Type:   mount.TypeBind,
							Source: path,
							Target: "/app",
						})

						// get additional site mounts
						siteMounts, err := site.GetAbsMountPaths(home)
						if err != nil {
							return err
						}

						// create mounts for the site
						for k, v := range siteMounts {
							mounts = append(mounts, mount.Mount{
								Type:   mount.TypeBind,
								Source: k,
								Target: v,
							})
						}

						// check if xdebug is enabled
						if site.Xdebug {
							envs = append(envs, "XDEBUG_MODE=develop,debug")
						} else {
							envs = append(envs, "XDEBUG_MODE=off")
						}

						// create new container, will have a new container id
						resp, err := docker.ContainerCreate(
							ctx,
							&container.Config{
								Image: image,
								Labels: map[string]string{
									labels.Environment: env,
									labels.Host:        site.Hostname,
								},
								Env: envs,
							},
							&container.HostConfig{
								Mounts: mounts,
							},
							&network.NetworkingConfig{
								EndpointsConfig: map[string]*network.EndpointSettings{
									env: {
										NetworkID: networkID,
									},
								},
							},
							site.Hostname,
						)
						if err != nil {
							return fmt.Errorf("unable to create the container, %w", err)
						}

						output.Pending("starting", site.Hostname)

						// start the container
						if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
							output.Warning()
							return fmt.Errorf("unable to start the container, %w", err)
						}

						output.Done()

						break
					}
				default:
					// create a brand new container since there is not an existing one
					image := fmt.Sprintf(NginxImage, site.PHP)

					// should we skip pulling the image
					if cmd.Flag("skip-pull").Value.String() == "false" {
						output.Pending("pulling", image)

						rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
						if err != nil {
							return fmt.Errorf("unable to pull the image, %w", err)
						}

						buf := &bytes.Buffer{}
						if _, err := buf.ReadFrom(rdr); err != nil {
							return fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
						}

						output.Done()
					}

					// get the sites main path
					path, err := site.GetAbsPath(home)
					if err != nil {
						return err
					}

					output.Pending("creating", site.Hostname)

					// add the path mount
					mounts := []mount.Mount{}
					mounts = append(mounts, mount.Mount{
						Type:   mount.TypeBind,
						Source: path,
						Target: "/app",
					})

					// get additional site mounts
					siteMounts, err := site.GetAbsMountPaths(home)
					if err != nil {
						return err
					}

					for k, v := range siteMounts {
						mounts = append(mounts, mount.Mount{
							Type:   mount.TypeBind,
							Source: k,
							Target: v,
						})
					}

					// check if xdebug is enabled
					if site.Xdebug {
						envs = append(envs, "XDEBUG_MODE=develop,debug")
					} else {
						envs = append(envs, "XDEBUG_MODE=off")
					}

					// create the container
					resp, err := docker.ContainerCreate(
						ctx,
						&container.Config{
							Image: image,
							Labels: map[string]string{
								labels.Environment: env,
								labels.Host:        site.Hostname,
							},
							Env: envs,
						},
						&container.HostConfig{
							Mounts: mounts,
						},
						&network.NetworkingConfig{
							EndpointsConfig: map[string]*network.EndpointSettings{
								env: {
									NetworkID: networkID,
								},
							},
						},
						site.Hostname,
					)
					if err != nil {
						return fmt.Errorf("unable to create the container, %w", err)
					}

					output.Pending("starting", site.Hostname)

					// start the container
					if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
						output.Warning()
						return fmt.Errorf("unable to start the container, %w", err)
					}

					output.Done()
				}

				// remove the site filter
				filter.Del("label", labels.Host+"="+site.Hostname)
			}

			// convert the sites into the gRPC API Apply request
			sites := make(map[string]*protob.Site)
			for _, s := range cfg.Sites {
				hosts := []string{s.Hostname}

				// if there are aliases lets append them to the hosts
				if len(s.Aliases) > 0 {
					hosts = append(hosts, s.Aliases...)
				}

				// create the site
				sites[s.Hostname] = &protob.Site{
					Hostname: s.Hostname,
					Aliases:  strings.Join(hosts, ","),
					Port:     8080,
				}
			}

			// ping the nitrod API until its ready...
			output.Pending("waiting for api")
			ping := &protob.PingRequest{}
			waiting := true
			for waiting {
				_, err := nitrod.Ping(ctx, ping)
				if err == nil {
					waiting = false
				}
			}

			output.Done()

			// configure the proxy with the sites
			output.Info("Configuring Proxy...")
			if _, err = nitrod.Apply(ctx, &protob.ApplyRequest{Sites: sites}); err != nil {
				return err
			}

			output.Success("proxy ready")

			// get all possible hostnames
			var hostnames []string
			for _, s := range cfg.Sites {
				hostnames = append(hostnames, s.Hostname)
				hostnames = append(hostnames, s.Aliases...)
			}

			// get the executable
			nitro, err := os.Executable()
			if err != nil {
				return err
			}

			// run the hosts command
			switch runtime.GOOS {
			case "windows":
				return fmt.Errorf("setting hosts file is not yet supported on windows")
			default:
				output.Info("Modifying hosts file (you might be prompted for your password)")

				// add the hosts
				if err := sudo.Run(nitro, "nitro", "hosts", "--hostnames="+strings.Join(hostnames, ",")); err != nil {
					return err
				}
			}

			output.Info(env, "is up and running 😃")

			return nil
		},
	}

	// add flag to skip pulling images
	cmd.Flags().BoolP("skip-pull", "s", false, "skip pulling images")

	return cmd
}
