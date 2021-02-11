package blackfire

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # setup blackfire credentials
  nitro blackfire`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "blackfire",
		Short:   "Setup Blackfire credentials",
		Example: exampleText,
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.RunApply(cmd, args, output)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.VerifyInit(cmd, args, home, output)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// load the configuration
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// ask for the
			id, err := output.Ask("Enter your Blackfire Server ID", "", ":", nil)
			if err != nil {
				return err
			}

			token, err := output.Ask("Enter your Blackfire Server Token", "", ":", nil)
			if err != nil {
				return err
			}

			cfg.Blackfire.ServerID = id
			cfg.Blackfire.ServerToken = token

			// save the config file
			if err := cfg.Save(); err != nil {
				return fmt.Errorf("unable to save config, %w", err)
			}

			return nil
		},
	}

	return cmd
}
