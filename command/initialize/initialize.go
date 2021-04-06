package initialize

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/proxycontainer"
	"github.com/craftcms/nitro/pkg/setup"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # setup nitro
  nitro init`

var skipApply, skipTrust bool

// NewCommand takes a docker client and returns the init command for creating a new environment
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "init",
		Short:         "Performs Nitro’s initial setup.",
		Example:       exampleText,
		SilenceErrors: false,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// is the docker api alive?
			if _, err := docker.Ping(cmd.Context()); err != nil {
				return fmt.Errorf("Couldn’t connect to Docker; please make sure Docker is running.")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			// check if there is a config file
			_, err := config.Load(home)
			if errors.Is(err, config.ErrNoConfigFile) {
				// walk the user through the first time setup
				if err := setup.FirstTime(home, cmd.InOrStdin(), output); err != nil {
					return err
				}
			}

			output.Info("Checking Nitro…")

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
						containerlabels.Nitro:   "true",
						containerlabels.Network: "true",
					},
				})
				if err != nil {
					return fmt.Errorf("unable to create the network, %w", err)
				}

				// set the newly created network
				networkID = resp.ID

				output.Done()
			}

			// create the proxy container
			if err := proxycontainer.Create(cmd.Context(), docker, output, networkID); err != nil {
				return err
			}

			// run the follow up commands
			for _, c := range cmd.Root().Commands() {
				// should we run the apply command
				if !skipApply {
					if c.Use == "apply" {
						if err := c.RunE(c, args); err != nil {
							return err
						}
					}
				}

				if !skipTrust {
					// should we run the trust command
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
	cmd.Flags().BoolVar(&skipApply, "skip-apply", false, "skip applying changes")
	cmd.Flags().BoolVar(&skipTrust, "skip-trust", false, "skip trusting the root certificate")

	return cmd
}
