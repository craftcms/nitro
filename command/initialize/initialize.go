package initialize

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/proxycontainer"
	"github.com/craftcms/nitro/pkg/setup"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # setup nitro
  nitro init`

// NewCommand takes a docker client and returns the init command for creating a new environment
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Setup nitro",
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

			// create the proxy container
			if err := proxycontainer.Create(cmd.Context(), docker, output, volume, networkID); err != nil {
				return err
			}

			// convert the apply flag to a boolean
			skipApply, err := strconv.ParseBool(cmd.Flag("skip-apply").Value.String())
			if err != nil {
				// don't do anything with the error
				skipApply = false
			}

			// check if we need to run the
			if !skipApply && cmd.Parent() != nil {
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
