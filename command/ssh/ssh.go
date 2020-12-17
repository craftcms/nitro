package ssh

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
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

// https://github.com/moby/moby/blob/8e610b2b55bfd1bfa9436ab110d311f5e8a74dcb/integration/internal/container/exec.go#L38

const exampleText = `  # ssh into a container - assuming its the current working directory
  nitro ssh`

// NewCommand returns the ssh command to get a shell in a container. The command is context aware and if
// it is not in a known project directory, it will provide a list of known sites to the user.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
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
			var site string
			var sites []string
			for _, s := range cfg.Sites {
				p, _ := s.GetAbsPath(home)

				// check if the path matches a sites path, then we are in a known site
				if strings.Contains(wd, p) {
					site = s.Hostname
					break
				}

				// add the site to the list in case we cannot find the directory
				sites = append(sites, s.Hostname)
			}

			// check the current site
			switch site == "" {
			case true:
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", sites)
				if err != nil {
					return err
				}

				// add the label to get the site
				filter.Add("label", labels.Host+"="+sites[selected])
			default:
				// add the label to get the site
				filter.Add("label", labels.Host+"="+site)
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

			// find the docker executable
			cli, err := exec.LookPath("docker")
			if err != nil {
				return err
			}

			c := exec.Command(cli, "exec", "-it", containers[0].ID, "sh")
			c.Stdin = cmd.InOrStdin()
			c.Stderr = cmd.ErrOrStderr()
			c.Stdout = cmd.OutOrStdout()

			if err := c.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func connect(ctx context.Context, docker client.ContainerAPIClient, containerID string) error {
	inout := make(chan []byte)
	errCh := make(chan error)

	// attach to the container
	waiter, err := docker.ContainerAttach(ctx, containerID, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return fmt.Errorf("error attaching to the container, %w", err)
	}

	if err := docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		errCh <- err
	}

	go func() {
		scanner := bufio.NewScanner(os.Stdin)

		_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, waiter.Reader)
		errCh <- err
		for scanner.Scan() {
			fmt.Println("scanner")
			inout <- []byte(scanner.Text())
		}
	}()

	// Write to docker container
	go func(w io.WriteCloser) {
		for {
			data, ok := <-inout
			if !ok {
				fmt.Println("!ok")
				w.Close()
				return
			}

			w.Write(append(data, '\n'))
		}
	}(waiter.Conn)

	if _, err := docker.ContainerWait(ctx, containerID); err != nil {
		errCh <- err
	}

	return nil
}
