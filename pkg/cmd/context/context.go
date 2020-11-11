package context

import (
	"fmt"

	"github.com/craftcms/nitro/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ContextCommand = &cobra.Command{
	Use:   "context",
	Short: "View the environment information",
	RunE:  contextMain,
	Example: `  # view all resources for the environment
  nitro context`,
}

func contextMain(cmd *cobra.Command, args []string) error {
	env := cmd.Flag("environment").Value
	cfg := config.Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("unable to read the config file, %w", err)
	}

	fmt.Println("Using configuration file")
	fmt.Println("  ==>", viper.ConfigFileUsed())
	fmt.Println("")

	fmt.Println("The follow sites are configured for", env)
	for _, site := range cfg.Sites {
		fmt.Println("  ==> hostname: ", site.Hostname)
		fmt.Println("  ==> webroot: ", site.Webroot)
		if len(site.Aliases) > 0 {
			fmt.Println("  ==> aliases: ", site.Aliases)
		}
		fmt.Println("")
	}

	return nil
}
