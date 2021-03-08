package apply

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command/apply/internal/customcontainer"
	"github.com/craftcms/nitro/command/apply/internal/databasecontainer"
	"github.com/craftcms/nitro/command/apply/internal/sitecontainer"
	"github.com/craftcms/nitro/pkg/backup"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/wsl"

	"github.com/craftcms/nitro/pkg/datetime"
	"github.com/craftcms/nitro/pkg/hostedit"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/proxycontainer"
	"github.com/craftcms/nitro/pkg/sudo"
	"github.com/craftcms/nitro/pkg/svc/dynamodb"
	"github.com/craftcms/nitro/pkg/svc/mailhog"
	"github.com/craftcms/nitro/pkg/svc/minio"
	"github.com/craftcms/nitro/pkg/svc/redis"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/craftcms/nitro/protob"
)

var (
	defaultFile     = "/etc/hosts"
	hostnames       []string
	knownContainers = map[string]bool{}
	isWSL           = false
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// is the docker api alive?
			if _, err := docker.Ping(cmd.Context()); err != nil {
				return fmt.Errorf("Couldn’t connect to Docker; please make sure Docker is running.")
			}

			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", labels.Nitro+"=true")

			// look for a container for the site
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{All: true, Filters: filter})
			if err != nil {
				return fmt.Errorf("error getting a list of containers")
			}

			if len(containers) > 0 {
				output.Info("Cleaning up...")
			}

			for _, c := range containers {
				// start the container if not running
				if c.State != "running" {
					for _, command := range cmd.Root().Commands() {
						if command.Use == "start" {
							if err := command.RunE(cmd, []string{}); err != nil {
								return err
							}
						}
					}
				}

				if _, ok := knownContainers[c.ID]; !ok {
					// don't remove the proxy container
					if c.Labels[labels.Proxy] != "" {
						continue
					}

					// set the container name
					name := strings.TrimLeft(c.Names[0], "/")

					output.Pending("removing", name)

					// only perform a backup if the container is for databases
					if c.Labels[labels.DatabaseEngine] != "" {
						// get all of the databases
						databases, err := backup.Databases(cmd.Context(), docker, c.ID, c.Labels[labels.DatabaseCompatibility])
						if err != nil {
							output.Warning()
							output.Info("Unable to get the databases from", name, err.Error())
							break
						}

						// backup each database
						for _, db := range databases {
							// create the database specific backup options
							opts := &backup.Options{
								BackupName:    fmt.Sprintf("%s-%s.sql", db, datetime.Parse(time.Now())),
								ContainerID:   c.ID,
								ContainerName: name,
								Database:      db,
								Home:          home,
							}

							// create the backup command based on the compatibility type
							switch c.Labels[labels.DatabaseCompatibility] {
							case "postgres":
								opts.Commands = []string{"pg_dump", "--username=nitro", db, "-f", "/tmp/" + opts.BackupName}
							default:
								opts.Commands = []string{"/usr/bin/mysqldump", "-h", "127.0.0.1", "-unitro", "--password=nitro", db, "--result-file=" + "/tmp/" + opts.BackupName}
							}

							output.Pending("creating backup", opts.BackupName)

							// backup the container
							if err := backup.Perform(cmd.Context(), docker, opts); err != nil {
								output.Warning()
								output.Info("Unable to backup database", db, err.Error())
								break
							}

							output.Done()
						}

						// show where all backups are saved for this container
						output.Info("Backups saved in", filepath.Join(home, config.DirectoryName, name), "💾")
					}

					// stop and remove a container we don't know about
					if err := docker.ContainerStop(cmd.Context(), c.ID, nil); err != nil {
						return err
					}

					// remove container
					if err := docker.ContainerRemove(cmd.Context(), c.ID, types.ContainerRemoveOptions{}); err != nil {
						return err
					}

					output.Done()
				}
			}

			if isWSL {
				output.Info(fmt.Sprintf("For your hostnames to work, add the following to `%s`:", `C:\Windows\System32\Drivers\etc\hosts`))
				output.Info("---- COPY BELOW ----")
				output.Info(fmt.Sprintf(`# <nitro>
%s %s
# </nitro>`, "127.0.0.1", strings.Join(hostnames, " ")))
				output.Info("---- COPY ABOVE ----")
			}

			output.Info("Nitro is up and running 😃")

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				// when we call commands from other commands (e.g. init)
				// the context could be nil, so we set it to the parent
				// context just in case.
				ctx = context.Background()
			}

			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", labels.Nitro+"=true")

			// add the filter for the network name
			filter.Add("name", "nitro-network")

			output.Info("Checking network…")

			// check the network
			var network types.NetworkResource
			networks, err := docker.NetworkList(ctx, types.NetworkListOptions{Filters: filter})
			if err != nil {
				return fmt.Errorf("unable to list docker networks\n%w", err)
			}

			// get the network for the environment
			for _, n := range networks {
				if n.Name == "nitro-network" {
					network = n
					break
				}
			}

			// if the network is not found
			if network.ID == "" {
				output.Info("No network was found…\nrun `nitro init` to get started")
				return nil
			}

			// remove the filter
			filter.Del("name", "nitro-network")

			output.Success("network ready")

			output.Info("Checking proxy…")

			// check the proxy and ensure its started
			_, err = proxycontainer.FindAndStart(ctx, docker)
			if errors.Is(err, proxycontainer.ErrNoProxyContainer) {
				// create the proxy
				if err := proxycontainer.Create(ctx, docker, output, network.ID); err != nil {
					output.Info("unable to find the nitro proxy…\n run `nitro init` to resolve")
					return err
				}
			}
			if err != nil && !errors.Is(err, proxycontainer.ErrNoProxyContainer) {
				return err
			}

			output.Success("proxy ready")

			output.Info("Checking databases…")

			// check the databases
			for _, db := range cfg.Databases {
				n, _ := db.GetHostname()
				output.Pending("checking", n)

				// start or create the database
				id, hostname, err := databasecontainer.StartOrCreate(ctx, docker, network.ID, db, output)
				if err != nil {
					output.Warning()
					return err
				}

				// set the container as known
				knownContainers[id] = true

				// add the hostname to the hosts files
				hostnames = append(hostnames, hostname)

				output.Done()
			}

			output.Info("Checking services…")

			// check dynamodb service
			switch cfg.Services.DynamoDB {
			case false:
				output.Pending("checking dynamodb service")

				if err := dynamodb.VerifyRemoved(ctx, docker, output); err != nil {
					output.Warning()
					return err
				}

				output.Done()
			default:
				output.Pending("checking dynamodb service")

				id, hostname, err := dynamodb.VerifyCreated(ctx, docker, network.ID, output)
				if err != nil {
					return err
				}

				if id != "" {
					knownContainers[id] = true
				}

				if hostname != "" {
					hostnames = append(hostnames, hostname)
				}

				output.Done()
			}

			// check mailhog service
			switch cfg.Services.Mailhog {
			case false:
				output.Pending("checking mailhog service")

				// make sure the service container is removed
				if err := mailhog.VerifyRemoved(ctx, docker, output); err != nil {
					return err
				}

				output.Done()
			default:
				output.Pending("checking mailhog service")

				// verify the mailhog container is created
				id, hostname, err := mailhog.VerifyCreated(ctx, docker, network.ID, output)
				if err != nil {
					return err
				}

				if id != "" {
					knownContainers[id] = true
				}

				if hostname != "" {
					hostnames = append(hostnames, hostname)
				}

				output.Done()
			}

			// check minio service
			switch cfg.Services.Minio {
			case false:
				// make sure the service container is removed
				err := minio.VerifyRemoved(ctx, docker, output)
				if err != nil {
					return err
				}
			default:
				output.Pending("checking minio service")

				// verify the minio container is created
				id, hostname, err := minio.VerifyCreated(ctx, docker, network.ID, output)
				if err != nil {
					return err
				}

				if id != "" {
					knownContainers[id] = true
				}

				if hostname != "" {
					hostnames = append(hostnames, hostname)
				}

				output.Done()
			}

			// check redis service
			switch cfg.Services.Redis {
			case false:
				output.Pending("checking redis service")

				if err := redis.VerifyRemoved(ctx, docker, output); err != nil {
					return err
				}

				output.Done()
			default:
				output.Pending("checking redis service")

				id, hostname, err := redis.VerifyCreated(ctx, docker, network.ID, output)
				if err != nil {
					return err
				}

				if id != "" {
					knownContainers[id] = true
				}

				if hostname != "" {
					hostnames = append(hostnames, hostname)
				}

				output.Done()
			}

			if len(cfg.Containers) > 0 {
				// get all of the containers
				output.Info("Checking containers...")

				for _, c := range cfg.Containers {
					output.Pending("checking", fmt.Sprintf("%s.containers.nitro", c.Name))

					// start, update or create the custom container
					id, err := customcontainer.StartOrCreate(ctx, docker, home, network.ID, c)
					if err != nil {
						output.Warning()
						return err
					}

					knownContainers[id] = true

					output.Done()
				}
			}

			if len(cfg.Sites) > 0 {
				// get all of the sites, their local path, the php version, and the type of project (nginx or PHP-FPM)
				output.Info("Checking sites…")

				// get the envs for the sites
				for _, site := range cfg.Sites {
					output.Pending("checking", site.Hostname)

					// start, update or create the site container
					id, err := sitecontainer.StartOrCreate(ctx, docker, home, network.ID, site, cfg)
					if err != nil {
						output.Warning()
						return err
					}

					knownContainers[id] = true

					output.Done()
				}
			}

			output.Info("Checking proxy…")

			output.Pending("updating proxy")

			if err := updateProxy(ctx, docker, nitrod, cfg); err != nil {
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
			for _, s := range cfg.Sites {
				hostnames = append(hostnames, s.Hostname)
				hostnames = append(hostnames, s.Aliases...)
			}

			// get custom container hostnames
			for _, c := range cfg.Containers {
				hostnames = append(hostnames, fmt.Sprintf("%s.containers.nitro", c.Name))
			}

			if len(hostnames) > 0 {
				// is this wsl?
				isWSL = wsl.IsWSL()

				// set the hosts file based on the OS
				if runtime.GOOS == "windows" {
					defaultFile = `C:\Windows\System32\Drivers\etc\hosts`
				}

				// check if hosts is already up to date
				updated, err := hostedit.IsUpdated(defaultFile, "127.0.0.1", hostnames...)
				if err != nil {
					return err
				}

				// if the hosts file is not updated
				if !updated {
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
						output.Info("Updating hosts file (you might be prompted for your password)")

						// add the hosts
						if err := sudo.Run(nitro, "nitro", "hosts", "--hostnames="+strings.Join(hostnames, ",")); err != nil {
							return err
						}
					}
				}
			}

			return nil
		},
	}

	// add flag to skip pulling images
	cmd.Flags().Bool("skip-hosts", false, "skip modifying the hosts file")

	return cmd
}

func updateProxy(ctx context.Context, docker client.ContainerAPIClient, nitrod protob.NitroClient, cfg *config.Config) error {
	// convert the sites into the gRPC API Apply request
	sites := make(map[string]*protob.Site)
	for _, s := range cfg.Sites {
		// create the site
		sites[s.Hostname] = &protob.Site{
			Hostname: s.Hostname,
			Aliases:  strings.Join(s.Aliases, ","),
			Port:     8080,
		}
	}

	// check the mailhog service
	if cfg.Services.Mailhog {
		sites["mailhog.service.nitro"] = &protob.Site{
			Hostname: "mailhog.service.nitro",
			Port:     8025,
		}
	}

	// check the minio service
	if cfg.Services.Minio {
		sites["minio.service.nitro"] = &protob.Site{
			Hostname: "minio.service.nitro",
			Port:     9000,
		}
	}

	// add any custom containers that need to be proxied
	for _, c := range cfg.Containers {
		if c.WebGui != 0 {
			sites[fmt.Sprintf("%s.containers.nitro", c.Name)] = &protob.Site{
				Hostname: fmt.Sprintf("%s.containers.nitro", c.Name),
				Port:     int32(c.WebGui),
			}
		}
	}

	// if there are no sites, we are done
	if len(sites) == 0 {
		return nil
	}

	// wait for the api to be ready
	for {
		_, err := nitrod.Ping(ctx, &protob.PingRequest{})
		if err == nil {
			break
		}
	}

	// configure the proxy with the sites
	resp, err := nitrod.Apply(ctx, &protob.ApplyRequest{Sites: sites})
	if err != nil {
		return err
	}

	if resp.Error {
		return fmt.Errorf("unable to update the proxy, %s", resp.GetMessage())
	}

	return nil
}
