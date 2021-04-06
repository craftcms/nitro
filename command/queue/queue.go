package queue

import (
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # execute the craft queue command for a site
  nitro queue`

// NewCommand returns the command to run queue listen inside of a sites container. It will check if the
// current working directory is a known site and auto-select or prompt a user for a list of sites.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "queue",
		Short:              "Runs a queue worker.",
		DisableFlagParsing: true,
		Example:            exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			// get a context aware list of sites
			sites := cfg.ListOfSitesByDirectory(home, wd)

			// create the options for the sites
			var options []string
			for _, s := range sites {
				options = append(options, s.Hostname)
			}

			// check if we found a site
			var site config.Site
			switch len(sites) {
			case 0:
				// prompt for the site
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
				if err != nil {
					return err
				}

				// set the site we selected
				site = sites[selected]

				// add the label to get the site
				filter.Add("label", containerlabels.Host+"="+sites[selected].Hostname)
			case 1:
				output.Info("connecting to", sites[0].Hostname)

				// set the site we selected
				site = sites[0]

				// add the label to get the site
				filter.Add("label", containerlabels.Host+"="+sites[0].Hostname)
			default:
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
				if err != nil {
					return err
				}

				// set the site we selected
				site = sites[selected]

				// add the label to get the site
				filter.Add("label", containerlabels.Host+"="+sites[selected].Hostname)
			}

			// find the containers but limited to the site label
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter, All: true})
			if err != nil {
				return err
			}

			// are there any containers??
			if len(containers) == 0 {
				return fmt.Errorf("unable to find an matching site")
			}

			// start the container if its not running
			if containers[0].State != "running" {
				for _, command := range cmd.Root().Commands() {
					if command.Use == "start" {
						if err := command.RunE(cmd, []string{}); err != nil {
							return err
						}
					}
				}
			}

			// get the container path
			var commands []string
			path := site.GetContainerPath()
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
