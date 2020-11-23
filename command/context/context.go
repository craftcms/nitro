package context

import (
	"fmt"
	"strings"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const exampleText = `  # view all resources for the environment
  nitro context`

// New is used for scaffolding new commands
func New(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "context",
		Short:   "View environment information",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Config{}
			if err := viper.Unmarshal(&cfg); err != nil {
				return fmt.Errorf("unable to read config file, %w", err)
			}

			output.Info("Craft Nitro", "2.0.0")
			output.Info("")
			output.Info("Configuration:\t", viper.ConfigFileUsed())
			output.Info("")

			output.Info(`Sites:`)
			for _, site := range cfg.Sites {
				output.Info("    hostname:\t", site.Hostname)
				if len(site.Aliases) > 0 {
					output.Info("    aliases:\t", strings.Join(site.Aliases, ", "))
				}
				output.Info("    php:\t", site.PHP)
				output.Info("    webroot:\t", site.Dir)
				output.Info("    path:\t", site.Path)
				output.Info("")
			}

			output.Info(`Databases:`)
			for _, db := range cfg.Databases {
				output.Info("    engine:\t", db.Engine, db.Version)
				output.Info("    username:\t", "nitro")
				output.Info("    password:\t", "nitro")
				output.Info("    port:\t", db.Port)
				output.Info("")
			}

			return nil
		},
	}

	return cmd
}
