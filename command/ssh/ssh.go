package ssh

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/labels"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var (
	// ErrExample is used when we want to share an error
	ErrExample = fmt.Errorf("some example error")
)

const exampleText = `  # example command
  nitro ssh`

func New(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ssh",
		Short:   "SSH into a container",
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

				// check if the path matche
				if strings.Contains(p, wd) {
					filter.Add("label", labels.Host+"="+s.Hostname)

					// find all of the containers, there should only be one though
					containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter})
					if err != nil {
						return err
					}

					// create an exec
					exec, err := docker.ContainerExecCreate(cmd.Context(), containers[0].ID, types.ExecConfig{
						Tty:          true,
						AttachStdin:  true,
						AttachStderr: true,
						AttachStdout: true,
						Detach:       false,
						Cmd:          []string{"sh"},
					})
					if err != nil {
						return err
					}

					// attach to the exec
					resp, err := docker.ContainerExecAttach(cmd.Context(), exec.ID, types.ExecConfig{
						Tty:          true,
						AttachStdin:  true,
						AttachStderr: true,
						AttachStdout: true,
						Detach:       false,
						Cmd:          []string{"sh"},
					})

					// // start the exec
					// if err := docker.ContainerExecStart(cmd.Context(), exec.ID, types.ExecStartCheck{Detach: false, Tty: true}); err != nil {
					// 	return fmt.Errorf("unable to start the container, %w", err)
					// }

					rdr := resp.Reader
					for {
						rdr.R
					}

					if _, err := io.Copy(os.Stdin, resp.Reader); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}

	return cmd
}
