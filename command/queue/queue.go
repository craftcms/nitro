package queue

import (
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
)

// https://github.com/moby/moby/blob/8e610b2b55bfd1bfa9436ab110d311f5e8a74dcb/integration/internal/container/exec.go#L38

const exampleText = `  # example command
  nitro queue`

// NewCommand returns the command to run queue listen inside of a sites container. It will check if the
// current working directory is a known site and auto-select or prompt a user for a list of sites.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "queue",
		Short:   "Run a queue worker",
		Example: exampleText,
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
			filter.Add("label", labels.Nitro)

			// get all of the sites
			var sites, found []string
			for _, s := range cfg.Sites {
				p, _ := s.GetAbsPath(home)

				// check if the path matches a sites path, then we are in a known site
				if strings.Contains(wd, p) {
					found = append(found, s.Hostname)
				}

				// add the site to the list in case we cannot find the directory
				sites = append(sites, s.Hostname)
			}

			// check if we found a site
			switch len(found) {
			case 0:
				// prompt for the site
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", sites)
				if err != nil {
					return err
				}

				// add the label to get the site
				filter.Add("label", labels.Host+"="+sites[selected])
			case 1:
				output.Info("connecting to", found[0])

				// add the label to get the site
				filter.Add("label", labels.Host+"="+found[0])
			default:
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", found)
				if err != nil {
					return err
				}

				// add the label to get the site
				filter.Add("label", labels.Host+"="+found[selected])
			}

			// find the containers but limited to the site label
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter})
			if err != nil {
				return err
			}

			// are there any containers??
			if len(containers) == 0 {
				return fmt.Errorf("unable to find an matching site")
			}

			// start the container if its not running
			if containers[0].State != "running" {
				if err := docker.ContainerStart(cmd.Context(), containers[0].ID, types.ContainerStartOptions{}); err != nil {
					return err
				}
			}

			// get the host name from the container
			hostname := strings.TrimLeft(containers[0].Names[0], "/")

			// if the webroot is not the default and has a path seperator, assume its a directory
			s, err := cfg.FindSiteByHostName(hostname)
			if err != nil {
				return err
			}

			var cmds []string
			if strings.Contains(s.Webroot, "/") {
				parts := strings.Split(s.Webroot, "/")

				var containerPath string
				if len(parts) >= 2 {
					containerPath = strings.Join(parts[:len(parts)-1], "/")
				} else {
					containerPath = parts[0]
				}

				cmds = []string{"php", fmt.Sprintf("%s/%s", containerPath, "craft"), "queue/listen", "--verbose"}
			} else {
				cmds = []string{"php", "craft", "queue/listen", "--verbose"}
			}

			output.Info("Listening for queue jobs…")

			// create an exec
			exec, err := docker.ContainerExecCreate(cmd.Context(), containers[0].ID, types.ExecConfig{
				AttachStderr: true,
				AttachStdout: true,
				Cmd:          cmds,
			})
			if err != nil {
				return err
			}

			// attach to the exec
			resp, err := docker.ContainerExecAttach(cmd.Context(), exec.ID, types.ExecStartCheck{
				// Tty:          true,
			})
			if err != nil {
				return err
			}
			defer resp.Close()

			outputDone := make(chan error)

			go func() {
				_, err := stdcopy.StdCopy(cmd.OutOrStdout(), cmd.OutOrStderr(), resp.Reader)
				outputDone <- err
			}()

			select {
			case err := <-outputDone:
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
