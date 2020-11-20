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

	fmt.Println("Craft Nitro", "2.0.0")
	fmt.Println("")
	fmt.Println("Configuration:\t", viper.ConfigFileUsed())
	fmt.Println("")

	fmt.Println(`Sites:`)
	for _, site := range cfg.Sites {
		fmt.Println("    hostname:\t", site.Hostname)
		if len(site.Aliases) > 0 {
			fmt.Println("    aliases:\t", site.Aliases)
		}
		fmt.Println("    php:\t", "7.4")
		fmt.Println("    webroot:\t", site.Webroot)
		fmt.Println("    local:\t", "~/dev/plugins-test.nitro")
		fmt.Println("")
	}

	fmt.Println(`Databases:`)
	for _, db := range cfg.Databases {
		fmt.Println("    engine:\t", db.Engine, db.Version)
		fmt.Println("    username:\t", "nitro")
		fmt.Println("    password:\t", "nitro")
		fmt.Println("    port:\t", db.Port)
		fmt.Println("")
	}

	return nil
}
