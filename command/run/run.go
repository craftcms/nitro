package run

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # run one off containers
  nitro run`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run",
		Short:   "Runs a container.",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
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

			// should the container bew removed after completion?
			if cmd.Flag("remove").Value.String() == "true" {
				c.Args = append(c.Args, "--rm")
			}

			// should the container be interactive
			if cmd.Flag("interactive").Value.String() == "true" {
				c.Args = append(c.Args, "-it")
			}

			// if the working dir is set, grab the current directory and mount it
			if cmd.Flag("working-dir").Value.String() != "" {
				// get the working dir
				wd, err := os.Getwd()
				if err != nil {
					return err
				}

				c.Args = append(c.Args, "-v")

				vol := fmt.Sprintf("%s:%s", wd, cmd.Flag("working-dir").Value.String())

				c.Args = append(c.Args, vol)
			}

			// set the image to use, if the image is not found docker will pull it
			c.Args = append(c.Args, cmd.Flag("image").Value.String())

			// append the args to the container
			c.Args = append(c.Args, args...)

			fmt.Println(c.Args)

			return c.Run()
		},
	}

	// set flags for the command
	cmd.Flags().String("working-dir", "", "sets the working directory for the container")
	cmd.Flags().Bool("interactive", true, "")
	cmd.Flags().String("image", "", "")
	cmd.Flags().Bool("remove", true, "")

	return cmd
}
