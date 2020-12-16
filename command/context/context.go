package context

import (
	"fmt"
	"strings"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

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
			// load the config file
			cfg, err := config.Load(home, env)
			if err != nil {
				return err
			}

			// if they are asking for yaml, show only the yaml
			if cmd.Flag("yaml").Value.String() == "true" {
				return yamlFmt(cfg)
			}

			output.Info("Craft Nitro", cmd.Parent().Version)
			output.Info("")
			output.Info("Configuration:\t", cfg.File)
			output.Info("")

			output.Info(`Sites:`)
			for _, site := range cfg.Sites {
				output.Info("  hostname:\t", site.Hostname)
				if len(site.Aliases) > 0 {
					output.Info("  aliases:\t", strings.Join(site.Aliases, ", "))
				}
				output.Info("  php:\t", site.PHP)
				output.Info("  webroot:\t", site.Dir)
				output.Info("  path:\t", site.Path)
				output.Info("  ---")
			}

			output.Info(`Databases:`)
			for _, db := range cfg.Databases {
				hostname, _ := db.GetHostname()
				output.Info("  engine:\t", db.Engine, db.Version, "\thostname:", hostname)
				output.Info("  username:\t", "nitro", "\tpassword:", "nitro")
				output.Info("  port:\t", db.Port)
				output.Info("  ---")
			}

			return nil
		},
	}

	cmd.Flags().Bool("yaml", false, "show the config file")

	return cmd
}

func yamlFmt(cfg *config.Config) error {
	// redact blackfire credentials
	if cfg.Blackfire.ServerID != "" {
		cfg.Blackfire.ServerID = "****************"
	}
	if cfg.Blackfire.ServerToken != "" {
		cfg.Blackfire.ServerToken = "********************************"
	}

	// marshal into the struct version so we can remove the blackfire credentials
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
