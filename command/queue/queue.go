package queue

import (
	"fmt"
	"os"

	"github.com/craftcms/nitro/pkg/appaware"
	"github.com/craftcms/nitro/pkg/appcontainer"
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

const exampleText = `  # execute the craft queue command for an app
  nitro queue`

// NewCommand returns the command to run queue listen inside an app container. It will check if the
// current working directory is a known app.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "queue",
		Short:              "Runs a queue worker.",
		DisableFlagParsing: true,
		Example:            exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

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

			output.Info("connecting to", appName)

			// add the label to get the app
			filter.Add("label", containerlabels.Host+"="+appName)

			// find the containers but limited to the app label
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter, All: true})
			if err != nil {
				return err
			}

			// are there any containers??
			if len(containers) == 0 {
				return fmt.Errorf("unable to find an matching app")
			}

			// start the container if it's not running
			if containers[0].State != "running" {
				for _, command := range cmd.Root().Commands() {
					if command.Use == "start" {
						if err := command.RunE(cmd, []string{}); err != nil {
							return err
						}
					}
				}
			}

			app, err := cfg.FindAppByHostname(appName)
			if err != nil {
				return err
			}

			// get the container path
			var commands []string
			path := appcontainer.ContainerPath(*app)
			if path != "" {
				commands = []string{"php", fmt.Sprintf("%s/%s", path, "craft"), "queue/listen", "--verbose"}
			} else {
				commands = []string{"php", "craft", "queue/listen", "--verbose"}
			}

			output.Info("Listening for queue jobs…")

			// create an exec
			exec, err := docker.ContainerExecCreate(cmd.Context(), containers[0].ID, types.ExecConfig{
				AttachStderr: true,
				AttachStdout: true,
				Cmd:          commands,
			})
			if err != nil {
				return err
			}

			// attach to the exec
			resp, err := docker.ContainerExecAttach(cmd.Context(), exec.ID, types.ExecStartCheck{})
			if err != nil {
				return err
			}
			defer resp.Close()

			done := make(chan error)
			go func() {
				_, err := stdcopy.StdCopy(cmd.OutOrStdout(), cmd.OutOrStderr(), resp.Reader)
				done <- err
			}()

			select {
			case err := <-done:
				if err != nil {
					return err
				}
				break
			case <-cmd.Context().Done():
				return cmd.Context().Err()
			}

			// get the exit code
			exit, err := docker.ContainerExecInspect(cmd.Context(), exec.ID)
			if err != nil {
				return err
			}

			// do something with the exit code
			output.Info(fmt.Sprintf("%d", exit.ExitCode))

			return nil
		},
	}

	return cmd
}
