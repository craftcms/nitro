package context

import (
	"strings"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/terminal"
)

const exampleText = `  # view all resources for the environment
  nitro context

  # show only the config file
  nitro context --yaml`

// New is used for scaffolding new commands
func New(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "context",
		Short:   "View environment information",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			env := cmd.Flag("environment").Value.String()
			cfg, err := config.Load(home, env)
			if err != nil {
				return err
			}

			output.Info("Craft Nitro", cmd.Parent().Version)
			output.Info("")
			output.Info("Configuration:\t", viper.ConfigFileUsed())
			output.Info("")

			if cmd.Flag("yaml").Value.String() == "true" {
				cfg := config.Config{}
				// marshal into the struct version so we can remove the blackfire credentials
				if err := viper.Unmarshal(&cfg); err != nil {
					return err
				}

				// redact blackfire credentials
				if cfg.Blackfire.ServerID != "" {
					cfg.Blackfire.ServerID = "****************"
				}
				if cfg.Blackfire.ServerToken != "" {
					cfg.Blackfire.ServerToken = "********************************"
				}

				out, err := yaml.Marshal(&cfg)
				if err != nil {
					return err
				}

				output.Info(string(out))

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

	cmd.Flags().BoolP("yaml", "p", false, "show the config file")

	return cmd
}
