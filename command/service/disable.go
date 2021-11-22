package service

import (
	"fmt"

	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/terminal"
)

// NewCommand returns the command to enable common nitro services. These services are provided as containers
// and do not require a user to configure the ports/volumes or images.
func disableCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disables a service.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				fmt.Println(cmd.UsageString())

				return fmt.Errorf("service name param missing")
			}

			return nil
		},
		ValidArgs: []string{"blackfire", "dynamodb", "mailhog", "minio", "redis"},
		Example: `  # disable services
  nitro service disable <service-name>

  # disable blackfire
  nitro service disable blackfire

  # disable mailhog
  nitro service disable mailhog

  # disable minio
  nitro service disable minio

  # disable dynamodb
  nitro service disable dynamodb`,
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.RunApply(cmd, args, false, output)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// load the configuration
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// disable the service
			switch args[0] {
			case "blackfire":
				cfg.Services.Blackfire = false
			case "dynamodb":
				cfg.Services.DynamoDB = false
			case "mailhog":
				cfg.Services.Mailhog = false
			case "minio":
				cfg.Services.Minio = false
			case "redis":
				cfg.Services.Redis = false
			default:
				return ErrUnknownService
			}

			// save the config file
			if err := cfg.Save(); err != nil {
				return fmt.Errorf("unable to save config, %w", err)
			}

			output.Info("Successfully disabled", args[0]+"!", "✨")

			return nil
		},
	}

	return cmd
}
