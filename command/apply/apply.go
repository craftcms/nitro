package apply

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command/apply/internal/databasecontainer"
	"github.com/craftcms/nitro/command/apply/internal/match"
	"github.com/craftcms/nitro/command/apply/internal/nginx"
	"github.com/craftcms/nitro/command/apply/internal/proxycontainer"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/hostedit"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/sudo"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/craftcms/nitro/protob"
)

var (
	// ErrNoNetwork is used when we cannot find the network
	ErrNoNetwork = fmt.Errorf("Unable to find the network")

	// ErrNoProxyContainer is returned when the proxy container is not found for an environment
	ErrNoProxyContainer = fmt.Errorf("unable to locate the proxy container")

	// NginxImage is the image used for sites, with the PHP version
	NginxImage = "docker.io/craftcms/nginx:%s-dev"

	// DatabaseImage is used for determining the engine and version
	DatabaseImage = "docker.io/library/%s:%s"
)

const exampleText = `  # apply changes from a config
  nitro apply

  # skip editing the hosts file
  nitro apply --skip-hosts

  # you can also set the environment variable "NITRO_EDIT_HOSTS" to "false"`

// NewCommand returns the command used to apply configuration file changes to a nitro environment.
func NewCommand(home string, docker client.CommonAPIClient, nitrod protob.NitroClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apply",
		Short:   "Apply changes",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				// when we call commands from other commands (e.g. init)
				// the context could be nil, so we set it to the parent
				// context just in case.
				ctx = cmd.Parent().Context()
			}

			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("name", "nitro-network")

			output.Info("Checking Network...")

			// check the network
			var envNetwork types.NetworkResource
			networks, err := docker.NetworkList(ctx, types.NetworkListOptions{})
			if err != nil {
				return fmt.Errorf("unable to list docker networks\n%w", err)
			}

			// get the network for the environment
			for _, n := range networks {
				if n.Name == "nitro-network" {
					envNetwork = n
					break
				}
			}

			// if the network is not found
			if envNetwork.ID == "" {
				return ErrNoNetwork
			}

			output.Success("network ready")

			output.Info("Checking Proxy...")

			// check the proxy and ensure its started
			if _, err := proxycontainer.FindAndStart(ctx, docker); err != nil {
				return err
			}

			output.Success("proxy ready")

			output.Info("Checking Databases...")

			// check the databases
			for _, db := range cfg.Databases {
				n, _ := db.GetHostname()
				output.Pending("checking", n)
				if err := databasecontainer.StartOrCreate(ctx, docker, envNetwork.ID, db); err != nil {
					output.Warning()
					return err
				}
				output.Done()
			}

			if len(cfg.Sites) > 0 {
				// get all of the sites, their local path, the php version, and the type of project (nginx or PHP-FPM)
				output.Info("Checking Sites...")

				// get the envs for the sites
				for _, site := range cfg.Sites {
					output.Pending("checking", site.Hostname)
					// TODO(jasonmccallister) check if this value should be set on linux hosts
					envs := site.AsEnvs("host.docker.internal")

					// add the site filter
					filter.Add("label", labels.Host+"="+site.Hostname)

					// look for a container for the site
					containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
					if err != nil {
						return fmt.Errorf("error getting a list of containers")
					}

					// if there are no containers we need to create one
					switch len(containers) == 0 {
					case true:
						// create the container
						image := fmt.Sprintf(NginxImage, site.Version)

						// pull the image
						rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
						if err != nil {
							return fmt.Errorf("unable to pull the image, %w", err)
						}

						buf := &bytes.Buffer{}
						if _, err := buf.ReadFrom(rdr); err != nil {
							return fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
						}

						// get the sites path
						path, err := site.GetAbsPath(home)
						if err != nil {
							return err
						}

						// add the site itself and any aliases to the extra hosts
						extraHosts := []string{fmt.Sprintf("%s:%s", site.Hostname, "127.0.0.1")}
						for _, s := range site.Aliases {
							extraHosts = append(extraHosts, fmt.Sprintf("%s:%s", s, "127.0.0.1"))
						}

						// create the container
						resp, err := docker.ContainerCreate(
							ctx,
							&container.Config{
								Image: image,
								Labels: map[string]string{
									labels.Host: site.Hostname,
								},
								Env: envs,
							},
							&container.HostConfig{
								Mounts: []mount.Mount{
									{
										Type:   mount.TypeBind,
										Source: path,
										Target: "/app",
									},
								},
								ExtraHosts: extraHosts,
							},
							&network.NetworkingConfig{
								EndpointsConfig: map[string]*network.EndpointSettings{
									"nitro": {
										NetworkID: envNetwork.ID,
									},
								},
							},
							nil,
							site.Hostname,
						)
						if err != nil {
							return fmt.Errorf("unable to create the container, %w", err)
						}

						// start the container
						if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
							return fmt.Errorf("unable to start the container, %w", err)
						}

						// TODO(jasonmccallister) check for a custom root and copt the template to the container
						if site.Dir != "web" {
							// create the nginx file
							conf := nginx.Generate(site.Dir)

							// create the temp file
							tr, err := archive.Generate("default.conf", conf)
							if err != nil {
								return err
							}

							// copy the file into the container
							if err := docker.CopyToContainer(ctx, resp.ID, "/tmp", tr, types.CopyToContainerOptions{AllowOverwriteDirWithFile: false}); err != nil {
								return err
							}

							commands := map[string][]string{
								"copy the file":       {"cp", "/tmp/default.conf", "/etc/nginx/conf.d/default.conf"},
								"set the permissions": {"chmod", "0644", "/etc/nginx/conf.d/default.conf"},
							}

							for _, c := range commands {
								// create the exec
								exec, err := docker.ContainerExecCreate(ctx, resp.ID, types.ExecConfig{
									User:         "root",
									AttachStdout: true,
									AttachStderr: true,
									Tty:          false,
									Cmd:          c,
								})
								if err != nil {
									return err
								}

								// attach to the container
								attach, err := docker.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{
									Tty: false,
								})
								defer attach.Close()

								// show the output to stdout and stderr
								if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, attach.Reader); err != nil {
									return fmt.Errorf("unable to copy the output of container, %w", err)
								}

								// start the exec
								if err := docker.ContainerExecStart(ctx, exec.ID, types.ExecStartCheck{}); err != nil {
									return fmt.Errorf("unable to start the container, %w", err)
								}

								// wait for the container exec to complete
								waiting := true
								for waiting {
									resp, err := docker.ContainerExecInspect(ctx, exec.ID)
									if err != nil {
										return err
									}

									waiting = resp.Running
								}
							}

							// start the container
							if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
								return fmt.Errorf("unable to start the container, %w", err)
							}
						}

						// remove the site filter
						filter.Del("label", labels.Host+"="+site.Hostname)
					default:
						// there is a running container
						c := containers[0]

						// get the containers details that include environment variables
						details, err := docker.ContainerInspect(ctx, c.ID)
						if err != nil {
							return err
						}

						// make sure container is in sync
						if match.Site(home, site, details) == false {
							fmt.Print("- updating... ")
							// stop container
							if err := docker.ContainerStop(ctx, c.ID, nil); err != nil {
								return err
							}

							// remove container
							if err := docker.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{}); err != nil {
								return err
							}

							// create the container
							image := fmt.Sprintf(NginxImage, site.Version)

							// pull the image
							rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
							if err != nil {
								return fmt.Errorf("unable to pull the image, %w", err)
							}

							buf := &bytes.Buffer{}
							if _, err := buf.ReadFrom(rdr); err != nil {
								return fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
							}

							// get the sites path
							path, err := site.GetAbsPath(home)
							if err != nil {
								return err
							}

							// add the site itself to the extra hosts
							extraHosts := []string{fmt.Sprintf("%s:%s", site.Hostname, "127.0.0.1")}
							for _, s := range site.Aliases {
								extraHosts = append(extraHosts, fmt.Sprintf("%s:%s", s, "127.0.0.1"))
							}

							// create the container
							resp, err := docker.ContainerCreate(
								ctx,
								&container.Config{
									Image: image,
									Labels: map[string]string{
										labels.Host: site.Hostname,
									},
									Env: envs,
								},
								&container.HostConfig{
									Mounts: []mount.Mount{
										{
											Type:   mount.TypeBind,
											Source: path,
											Target: "/app",
										},
									},
									ExtraHosts: extraHosts,
								},
								&network.NetworkingConfig{
									EndpointsConfig: map[string]*network.EndpointSettings{
										"nitro-network": {
											NetworkID: envNetwork.ID,
										},
									},
								},
								nil,
								site.Hostname,
							)
							if err != nil {
								return fmt.Errorf("unable to create the container, %w", err)
							}

							// start the container
							if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
								return fmt.Errorf("unable to start the container, %w", err)
							}
						}

						// remove the site filter
						filter.Del("label", labels.Host+"="+site.Hostname)
					}

					output.Done()
				}
			}

			output.Info("Checking Proxy...")

			output.Pending("updating proxy")

			if err := updateProxy(ctx, docker, nitrod, *cfg); err != nil {
				output.Warning()
				return err
			}

			output.Done()

			// should we update the hosts file?
			if os.Getenv("NITRO_EDIT_HOSTS") == "false" || cmd.Flag("skip-hosts").Value.String() == "true" {
				// skip updating the hosts file
				return nil
			}

			// get all possible hostnames
			var hostnames []string
			for _, s := range cfg.Sites {
				hostnames = append(hostnames, s.Hostname)
				hostnames = append(hostnames, s.Aliases...)
			}

			if len(hostnames) > 0 {
				// set the hosts file based on the OS
				defaultFile := "/etc/hosts"
				if runtime.GOOS == "windows" {
					defaultFile = `C:\Windows\System32\Drivers\etc\hosts`
				}

				// check if hosts is already up to date
				updated, err := hostedit.IsUpdated(defaultFile, "127.0.0.1", hostnames...)
				if err != nil {
					return err
				}

				// if the hosts file is not updated
				if updated == false {
					// get the executable
					nitro, err := os.Executable()
					if err != nil {
						return fmt.Errorf("unable to locate the nitro path, %w", err)
					}

					// run the hosts command
					switch runtime.GOOS {
					case "windows":
						// windows users should be running as admin, so just execute the hosts command
						// as is
						c := exec.Command(nitro, "hosts", "--hostnames="+strings.Join(hostnames, ","))

						c.Stdout = os.Stdout
						c.Stderr = os.Stderr

						if c.Run() != nil {
							return err
						}
					default:
						output.Info("Modifying hosts file (you might be prompted for your password)")

						// add the hosts
						if err := sudo.Run(nitro, "nitro", "hosts", "--hostnames="+strings.Join(hostnames, ",")); err != nil {
							return err
						}
					}
				}
			}

			output.Info("Nitro is up and running 😃")

			return nil
		},
	}

	// add flag to skip pulling images
	cmd.Flags().Bool("skip-hosts", false, "skip modifying the hosts file")

	return cmd
}

func checkProxy(ctx context.Context, docker client.ContainerAPIClient) (types.Container, error) {
	f := filters.NewArgs()
	f.Add("label", labels.Nitro)
	f.Add("label", labels.Type+"=proxy")
	// TODO(jasonmccallister) add the type filter as well?

	// check if there is an existing container for the nitro-proxy
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: f, All: true})
	if err != nil {
		return types.Container{}, fmt.Errorf("unable to list the containers\n%w", err)
	}

	// get the container id and determine if the container needs to start
	for _, c := range containers {
		for _, n := range c.Names {
			if n == "nitro-proxy" || n == "/nitro-proxy" {
				// check if it is running
				if c.State != "running" {
					if err := docker.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
						return types.Container{}, fmt.Errorf("unable to start the nitro container, %w", err)
					}
				}

				// return the container
				return c, nil
			}
		}
	}

	return types.Container{}, ErrNoProxyContainer
}

func checkDatabase(ctx context.Context, docker client.CommonAPIClient, output terminal.Outputer, db config.Database, networkID string) error {
	// add filters to check for the container
	filter := filters.NewArgs()
	filter.Add("label", labels.Nitro)
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
	switch len(containers) {
	// the database container exists
	case 1:
		// check if the container is running
		if containers[0].State != "running" {
			output.Pending("starting", hostname)

			// start the container
			if err := docker.ContainerStart(ctx, containers[0].ID, types.ContainerStartOptions{}); err != nil {
				output.Warning()
				return err
			}

			output.Done()
		} else {
			output.Success(hostname, "ready")
		}
	default:
		// database container does not exist, so create the volume and start it
		output.Pending("creating volume", hostname)

		// create the database labels
		lbls := map[string]string{
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
		volume, err := docker.VolumeCreate(ctx, volumetypes.VolumeCreateBody{Driver: "local", Name: hostname, Labels: lbls})
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
		output.Pending("pulling", image)

		// pull the image
		rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
		if err != nil {
			output.Warning()
			return fmt.Errorf("unable to pull image %s, %w", image, err)
		}

		// read the output to pull the image
		buf := &bytes.Buffer{}
		if _, err := buf.ReadFrom(rdr); err != nil {
			output.Warning()
			return fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
		}

		output.Done()

		output.Pending("creating", hostname)

		// set the port for the database
		port, err := nat.NewPort("tcp", db.Port)
		if err != nil {
			output.Warning()
			return fmt.Errorf("unable to create the port, %w", err)
		}
		containerConfig := &container.Config{
			Image:  image,
			Labels: lbls,
			ExposedPorts: nat.PortSet{
				port: struct{}{},
			},
			Env: envs,
		}
		hostConfig := &container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: volume.Name,
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
		}
		networkConfig := &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"nitro": {
					NetworkID: networkID,
				},
			},
		}

		// create the container for the database
		resp, err := docker.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, nil, hostname)
		if err != nil {
			return fmt.Errorf("unable to create the container, %w", err)
		}

		// start the container
		if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			return fmt.Errorf("unable to start the container, %w", err)
		}

		output.Done()
	}

	return nil
}

func updateProxy(ctx context.Context, docker client.ContainerAPIClient, nitrod protob.NitroClient, cfg config.Config) error {
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

	// if there are no sites, we are done
	if len(sites) == 0 {
		return nil
	}

	// wait for the api to be ready
	wait := true
	for wait {
		_, err := nitrod.Ping(ctx, &protob.PingRequest{})
		if err == nil {
			wait = false
			break
		}
	}

	// configure the proxy with the sites
	resp, err := nitrod.Apply(ctx, &protob.ApplyRequest{Sites: sites})
	if err != nil {
		return err
	}

	if resp.Error == true {
		return fmt.Errorf("unable to update the proxy, %s", resp.GetMessage())
	}

	return nil
}
