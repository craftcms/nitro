package apply

import (
	"fmt"

	"github.com/craftcms/nitro/internal/config"
	"github.com/craftcms/nitro/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ApplyCommand = &cobra.Command{
	Use:   "apply",
	Short: "Apply changes to an environment",
	RunE:  applyMain,
	Example: `  # apply changes from a config to the environment
  nitro apply`,
}

func applyMain(cmd *cobra.Command, args []string) error {
	env := cmd.Flag("environment").Value.String()
	cfg := config.Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("unable to read the config file, %w", err)
	}

	// create the new client
	nitro, err := client.NewClient()
	if err != nil {
		return fmt.Errorf("unable to create a client for docker, %w", err)
	}

	return nitro.Apply(cmd.Context(), env, cfg)
}
