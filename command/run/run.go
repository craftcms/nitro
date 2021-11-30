package run

import (
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
			path, err := exec.LookPath("docker")
			if err != nil {
				return err
			}

			c := exec.Command(path, "run")

			c.Stdin = cmd.InOrStdin()
			c.Stderr = cmd.ErrOrStderr()
			c.Stdout = cmd.OutOrStdout()

			if cmd.Flag("remove").Value.String() == "true" {
				c.Args = append(c.Args, "--rm")
			}

			if cmd.Flag("interactive").Value.String() == "true" {
				c.Args = append(c.Args, "-it")
			}

			c.Args = append(c.Args, cmd.Flag("image").Value.String())

			c.Args = append(c.Args, args...)

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
