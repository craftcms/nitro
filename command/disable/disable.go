package disable

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	// ErrUnknownService is used when an unknown service is requested
	ErrUnknownService = fmt.Errorf("unknown service requested")
)

const exampleText = `  # disable services
  nitro disable <service-name>

  # disable mailhog
  nitro disable mailhog

  # disable minio
  nitro disable minio

  # disable dynamodb
  nitro disable dynamodb`

// NewCommand returns the command to enable common nitro services. These services are provided as containers
// and do not require a user to configure the ports/volumes or images.
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
		ValidArgs: []string{"dynamodb", "mailhog", "minio", "redis"},
		Example:   exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// load the configuration
			cfg, err := config.Load(home, false)
			if err != nil {
				return err
			}

			// disable the service
			switch args[0] {
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
