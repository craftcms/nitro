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

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
)

// https://github.com/moby/moby/blob/8e610b2b55bfd1bfa9436ab110d311f5e8a74dcb/integration/internal/container/exec.go#L38

const exampleText = `  # example command
  nitro queue`

// NewCommand returns the command to run queue listen inside of a sites container. It will check if the
// current working directory is a known site and autoselect or prompt a user for a list of sites.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "queue",
		Short:   "Run a queue worker",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			env := cmd.Flag("environment").Value.String()

			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			cfg, err := config.Load(home, env)
			if err != nil {
				return err
			}

			// create a filter for the enviroment
			filter := filters.NewArgs()
			filter.Add("label", labels.Environment+"="+env)

			// get all of the sites
			for _, s := range cfg.Sites {
				p, _ := s.GetAbsPath(home)

				// check if the path matches a sites path, then we are in a known site
				if strings.Contains(wd, p) {
					output.Info("Listening for queue jobs...")

					filter.Add("label", labels.Host+"="+s.Hostname)

					// find all of the containers, there should only be one though
					containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter})
					if err != nil {
						return err
					}

					// create an exec
					exec, err := docker.ContainerExecCreate(cmd.Context(), containers[0].ID, types.ExecConfig{
						// AttachStdin:  true,
						AttachStderr: true,
						AttachStdout: true,
						Cmd:          []string{"./craft", "queue/listen", "--verbose"},
						// Tty:          true,
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
				}
			}

			return nil
		},
	}

	return cmd
}
