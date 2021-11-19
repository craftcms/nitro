package logs

import (
	"os"
	"strconv"

	"github.com/craftcms/nitro/pkg/appaware"
	"github.com/craftcms/nitro/pkg/flags"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # show logs from a site
  nitro logs

  # show only the last 5 minutes
  nitro logs --since 5m

  # show logs but don't follow
  nitro logs --follow=false`

// NewCommand returns the command to show a containers logs. It will check if the current working
// directory is a known app and default to that container or provide the user with a list of sites
// to view logs from. There are helpful flags such as since, timestamps, and follow that align with
// the docker logs API flags.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logs",
		Short:   "Displays container logs.",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			// get the app
			appName := flags.AppName
			if appName == "" {
				// get the current working directory
				wd, err := os.Getwd()
				if err != nil {
					return err
				}

				appName, err = appaware.Detect(*cfg, wd)
				if err != nil {
					return err
				}
			}

			output.Info("show logs for", appName)

			filter.Add("label", containerlabels.Host+"="+appName)

			// find all containers, there should only be one if we are in a known directory
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter})
			if err != nil {
				return err
			}

			// set the options for logging based on the command flags
			opts := types.ContainerLogsOptions{
				ShowStdout: true,
				ShowStderr: true,
			}

			// parse the flags
			timestamps, err := strconv.ParseBool(cmd.Flag("timestamps").Value.String())
			if err != nil {
				timestamps = false
			}
			opts.Timestamps = timestamps

			follow, err := strconv.ParseBool(cmd.Flag("follow").Value.String())
			if err != nil {
				follow = true
			}
			opts.Follow = follow

			if cmd.Flag("since").Value.String() != "" {
				opts.Since = cmd.Flag("since").Value.String()
			}

			// get the containers logs
			out, err := docker.ContainerLogs(cmd.Context(), containers[0].ID, opts)
			if err != nil {
				return err
			}

			// show the output
			if _, err := stdcopy.StdCopy(cmd.OutOrStdout(), cmd.ErrOrStderr(), out); err != nil {
				return err
			}

			return nil
		},
	}

	// set flags for the command
	cmd.Flags().Bool("follow", true, "follow log output")
	cmd.Flags().Bool("timestamps", false, "show timestamps")
	cmd.Flags().String("since", "", "Show logs since timestamp (e.g. 2013-01-02T13:23:37Z) or relative (e.g. 42m for 42 minutes)")

	return cmd
}
