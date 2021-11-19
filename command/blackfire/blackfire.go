package blackfire

import (
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # enable blackfire for an app
  nitro blackfire on

  # disable blackfire for an app
  nitro blackfire off`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "blackfire",
		Short:   "Manages Blackfire credentials.",
		Example: exampleText,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.VerifyInit(cmd, args, home, output)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(onCommand(home, docker, output), offCommand(home, docker, output))

	return cmd
}
