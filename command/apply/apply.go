package apply

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command/apply/internal/databasecontainer"
	"github.com/craftcms/nitro/command/apply/internal/proxycontainer"
	"github.com/craftcms/nitro/command/apply/internal/sitecontainer"
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
			filter.Add("label", labels.Nitro+"=true")

			// add the filter for the network name
			filter.Add("name", "nitro-network")

			output.Info("Checking Network...")

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
				output.Info("No network was found...")
				output.Info("run `nitro init` to get started")
				return nil
			}

			// remove the filter
			filter.Del("name", "nitro-network")

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

				// start or create the database
				if err := databasecontainer.StartOrCreate(ctx, docker, network.ID, db); err != nil {
					output.Warning()
					return err
				}

				output.Done()
			}

			// check dynamodb service
			if cfg.Services.DynamoDB {
				output.Pending("checking dynamodb service")

				if _, err := dynamodb(ctx, docker, output, cfg.Services.DynamoDB, network.ID); err != nil {
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

					// start, update or create the site container
					if err := sitecontainer.StartOrCreate(ctx, docker, home, network.ID, site); err != nil {
						output.Warning()
						return err
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

func updateProxy(ctx context.Context, docker client.ContainerAPIClient, nitrod protob.NitroClient, cfg config.Config) error {
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

	if resp.Error {
		return fmt.Errorf("unable to update the proxy, %s", resp.GetMessage())
	}

	return nil
}
