package run

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	flagImage, flagWorkingDir              string
	flagInteractive, flagPull, flagPersist bool
)

const exampleText = `  # run one off containers
  nitro run --image node:10 --working-dir /app install

  # run a composer container, mounting the current directory with a shell inside the container
  nitro run --image composer --working-dir /app bash`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run",
		Short:   "Runs a container.",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			filter := filters.NewArgs()
			filter.Add("reference", flagImage)

			// look for the image
			list, err := docker.ImageList(cmd.Context(), types.ImageListOptions{
				All:     true,
				Filters: filter,
			})
			if err != nil {
				return err
			}

			// should we pull the image?
			if len(list) == 0 || flagPull {
				fmt.Print("pulling image … ")

				// pull the image
				r, err := docker.ImagePull(cmd.Context(), flagImage, types.ImagePullOptions{})
				if err != nil {
					fmt.Print("\u2717\n")
					return err
				}
				defer r.Close()

				buf := bytes.Buffer{}
				if _, err := buf.ReadFrom(r); err != nil {
					fmt.Print("\u2717\n")
					return err
				}

				fmt.Print("\u2713\n")
			}

			// find the docker executable
			path, err := exec.LookPath("docker")
			if err != nil {
				return err
			}

			c := exec.Command(path, "run")

			// set stdout/stdin/stderr
			c.Stdin = cmd.InOrStdin()
			c.Stderr = cmd.ErrOrStderr()
			c.Stdout = cmd.OutOrStdout()

			// should the container be removed after completion?
			if flagPersist {
				c.Args = append(c.Args, "--rm")
			}

			// should the container be interactive
			if flagInteractive {
				c.Args = append(c.Args, "-it")
			}

			// if the working dir is set, grab the current directory and mount it
			if flagWorkingDir != "" {
				// get the working dir
				current, err := os.Getwd()
				if err != nil {
					return err
				}

				c.Args = append(c.Args, "-v")

				vol := fmt.Sprintf("%s:%s", current, flagWorkingDir)

				c.Args = append(c.Args, vol)
			}

			// set the image to use, if the image is not found docker will pull it
			c.Args = append(c.Args, flagImage)

			// append the args to the container
			c.Args = append(c.Args, args...)

			return c.Run()
		},
	}

	// set flags for the command
	cmd.Flags().StringVar(&flagWorkingDir, "working-dir", "", "working directory for the container")
	cmd.Flags().BoolVar(&flagInteractive, "interactive", true, "should the container be interactive?")
	cmd.Flags().StringVar(&flagImage, "image", "", "image to use for the container")
	cmd.Flags().BoolVar(&flagPersist, "persist", true, "persist container after completion")
	cmd.Flags().BoolVar(&flagPull, "pull", false, "pull the image, even if its been downloaded once")

	cmd.MarkFlagRequired("image")

	return cmd
}
