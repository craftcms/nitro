package ssh

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/labels"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/cobra"
)

var (
	// ErrExample is used when we want to share an error
	ErrExample = fmt.Errorf("some example error")
)

// https://github.com/moby/moby/blob/8e610b2b55bfd1bfa9436ab110d311f5e8a74dcb/integration/internal/container/exec.go#L38

const exampleText = `  # ssh into a container - assuming its the current working directory
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

			// load the config
			cfg, err := config.Load(home, env)
			if err != nil {
				return err
			}

			// create a filter for the enviroment
			filter := filters.NewArgs()
			filter.Add("label", labels.Environment+"="+env)

			// get all of the sites
			var containerID string
			var sites []string
			for _, s := range cfg.Sites {
				// add the site to the list in case we cannot find the directory
				sites = append(sites, s.Hostname)

				p, _ := s.GetAbsPath(home)

				// check if the path matches a sites path, then we are in a known site
				if strings.Contains(wd, p) {
					filter.Add("label", labels.Host+"="+s.Hostname)

					// find the containers but limited to the site label
					containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter})
					if err != nil {
						return err
					}

					// set the first container
					if len(containers) > 0 {
						containerID = containers[0].ID
					}
				}
			}

			// if we did not find a container, get a list of sites and prompt
			if containerID == "" {
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", sites)
				if err != nil {
					return err
				}

				// add the label to get the site
				filter.Add("label", labels.Host+"="+sites[selected])

				// find the containers but limited to the site label
				containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter})
				if err != nil {
					return err
				}

				// set the first container
				if len(containers) > 0 {
					containerID = containers[0].ID
				}
			}

			// connect to the container
			if err := connect(cmd.Context(), docker, containerID); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func connect(ctx context.Context, docker client.ContainerAPIClient, containerID string) error {
	// create an exec
	exec, err := docker.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          []string{"sh"},
		// Tty:          true,
	})
	if err != nil {
		return err
	}

	// attach to the exec
	stream, err := docker.ContainerExecAttach(ctx, exec.ID, types.ExecConfig{
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          []string{"sh"},
		// Tty:          true,
	})
	if err != nil {
		return err
	}
	defer stream.Close()

	outputDone := make(chan error)
	go func() {
		_, err := stdcopy.StdCopy(os.Stdout, os.Stderr, stream.Reader)
		outputDone <- err
	}()

	select {
	case err := <-outputDone:
		if err != nil {
			return err
		}
		break

	case <-ctx.Done():
		return ctx.Err()
	}

	// get the exit code
	exit, err := docker.ContainerExecInspect(ctx, exec.ID)
	if err != nil {
		return err
	}

	fmt.Println(exit)

	return nil
}
