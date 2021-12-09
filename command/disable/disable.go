package disable

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

// NewCommand returns the command to disable a service.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
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
  nitro disable <service-name>

  # disable blackfire
  nitro disable blackfire

  # disable mailhog
  nitro disable mailhog

  # disable minio
  nitro disable minio

  # disable dynamodb
  nitro disable dynamodb`,
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

			output.Info("Disabled", args[0])

			return nil
		},
	}

	return cmd
}
