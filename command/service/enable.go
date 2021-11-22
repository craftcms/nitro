package service

import (
	"fmt"

	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	// ErrUnknownService is used when an unknown service is requested
	ErrUnknownService = fmt.Errorf("unknown service requested")
)

func enableCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enables a service.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				fmt.Println(cmd.UsageString())

				return fmt.Errorf("service name param missing")
			}

			return nil
		},
		ValidArgs: []string{"blackfire", "dynamodb", "mailhog", "minio", "redis"},
		Example:   `  # enable services
  nitro service enable <service-name>

  # enable mailhog for local email testing
  nitro service enable mailhog

  # enable blackfire for local profiling
  nitro service enable blackfire

  # enable minio for local s3 testing
  nitro service enable minio

  # enable dynamodb for local noSQL
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

			// enable the service
			switch args[0] {
			case "blackfire":
				// TODO(jasonmccallister) verify the credentials are set
				cfg.Services.Blackfire = true
			case "dynamodb":
				cfg.Services.DynamoDB = true
			case "mailhog":
				cfg.Services.Mailhog = true
			case "minio":
				cfg.Services.Minio = true
			case "redis":
				cfg.Services.Redis = true
			default:
				return ErrUnknownService
			}

			// save the config file
			if err := cfg.Save(); err != nil {
				return fmt.Errorf("unable to save config, %w", err)
			}

			output.Info("Successfully enabled", args[0] + "!", "💥")

			return nil
		},
	}

	return cmd
}
