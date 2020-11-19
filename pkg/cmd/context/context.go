package context

import (
	"fmt"

	"github.com/craftcms/nitro/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ContextCommand is used to show contextual information
// based on the nitro environment such as sites, databases,
// and connection information
var ContextCommand = &cobra.Command{
	Use:   "context",
	Short: "View environment information",
	RunE:  contextMain,
	Example: `  # view all resources for the environment
  nitro context`,
}

func contextMain(cmd *cobra.Command, args []string) error {
	cfg := config.Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("unable to read config file, %w", err)
	}

	fmt.Println("Configuration")
	fmt.Println("  ==>", viper.ConfigFileUsed())
	fmt.Println("")

	fmt.Println("Databases")
	for _, db := range cfg.Databases {
		fmt.Println("  ==> engine:", db.Engine, "\tversion:", db.Version, "\tport:", db.Port)
		fmt.Println("      username:", "nitro", "\tpassword:", "nitro")
	}
	fmt.Println("")
	fmt.Println("Sites")
	for _, site := range cfg.Sites {
		// TODO(jasonmccallister) get the container information? Is it needed?
		fmt.Println("  ==> hostname: ", site.Hostname)
		if len(site.Aliases) > 0 {
			fmt.Println("      aliases: ", site.Aliases)
		}
		fmt.Println("      php:", "7.4", "\twebroot:", site.Webroot)

	}

	return nil
}
