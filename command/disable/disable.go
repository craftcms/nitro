package disable

import (
	"fmt"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var (
	// ErrUnknownService is used when an unknown service is requested
	ErrUnknownService = fmt.Errorf("unknown service requested")
)

const exampleText = `  # disable services
  nitro disable <service-name>

  # disable mailhog
  nitro disable mailhog

  # disable blackfire
  nitro disable blackfire

  # disable minio
  nitro disable minio

  # disable dynamodb
  nitro disable dynamodb`

// NewCommand returns the command to enable common nitro services. These services are provided as containers
// and do not require a user to configure the ports/volumes or images.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "disable",
		Short:     "Disable services",
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: []string{"blackfire", "dynamodb", "mailhog", "minio", "redis"},
		Example:   exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			env := cmd.Flag("environment").Value.String()

			// load the configuration
			cfg, err := config.Load(home, env)
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

			// run the apply command
			for _, c := range cmd.Parent().Commands() {
				// set the apply command
				if c.Use == "apply" {
					if err := c.RunE(c, args); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}

	return cmd
}
