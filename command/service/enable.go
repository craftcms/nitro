package service

import (
	"fmt"

	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/terminal"
)

// NewCommand returns the command to enable an app from automatically starting.
func enableCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enables an app.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				fmt.Println(cmd.UsageString())

				return fmt.Errorf("service name param missing")
			}

			return nil
		},
		ValidArgs: []string{"blackfire", "dynamodb", "mailhog", "minio", "redis"},
		Example: `  # enable services
  nitro service enable <service-name>

  # enable blackfire
  nitro service enable blackfire

  # enable mailhog
  nitro service enable mailhog

  # enable minio
  nitro service enable minio

  # enable dynamodb
  nitro service enable dynamodb`,
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

			output.Info("Enabled", args[0])

			return nil
		},
	}

	return cmd
}
