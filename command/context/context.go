package context

import (
	"io/ioutil"
	"strings"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const exampleText = `  # view all resources for the environment
  nitro context

  # show only the config file
  nitro context --not-pretty`

// New is used for scaffolding new commands
func New(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "context",
		Short:   "View environment information",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, cfg, err := config.Load()
			if err != nil {
				return err
			}

			output.Info("Craft Nitro", cmd.Parent().Version)
			output.Info("")
			output.Info("Configuration:\t", viper.ConfigFileUsed())
			output.Info("")

			if cmd.Flag("not-pretty").Value.String() == "true" {
				bytes, err := ioutil.ReadFile(viper.ConfigFileUsed())
				if err != nil {
					return err
				}

				output.Info(string(bytes))

				return nil
			}

			output.Info(`Sites:`)
			for _, site := range cfg.Sites {
				output.Info("    hostname:\t", site.Hostname)
				if len(site.Aliases) > 0 {
					output.Info("    aliases:\t", strings.Join(site.Aliases, ", "))
				}
				output.Info("    php:\t", site.PHP)
				output.Info("    webroot:\t", site.Dir)
				output.Info("    path:\t", site.Path)
				if len(site.Mounts) > 0 {
					output.Info("    mounts:\t", strings.Join(site.Mounts, ", "))
				}
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

	cmd.Flags().BoolP("not-pretty", "p", false, "show the not pretty version")

	return cmd
}
