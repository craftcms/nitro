package enable

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	// ErrUnknownService is used when an unknown service is requested
	ErrUnknownService = fmt.Errorf("unknown service requested")
)

const exampleText = `  # enable services
  nitro enable <service-name>

  # enable mailhog for local email testing
  nitro enable mailhog

  # enable blackfire for local profiling
  nitro enable blackfire

  # enable minio for local s3 testing
  nitro enable minio

  # enable dynamodb for local noSQL
  nitro enable dynamodb`

// NewCommand returns the command to enable common nitro services. These services are provided as containers
// and do not require a user to configure the ports/volumes or images.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "enable",
		Short:     "Enable services",
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

			// enable the service
			switch args[0] {
			case "blackfire":
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
