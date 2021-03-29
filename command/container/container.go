package container

import (
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # manage custom containers
  nitro container

  # add a new container
  nitro container new

  # ssh into a custom container
  nitro container ssh`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "container",
		Short:   "Manages custom containers.",
		Example: exampleText,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.VerifyInit(cmd, args, home, output)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		newCommand(home, docker, output),
		sshCommand(home, docker, output),
		removeCommand(home, docker, output),
	)

	return cmd
}
