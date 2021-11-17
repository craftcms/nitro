package extensions

import (
	"fmt"
	"os"
	"strings"

	"github.com/craftcms/nitro/pkg/appaware"
	"github.com/craftcms/nitro/pkg/flags"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # enable PHP extensions for the current app
  nitro extensions
  
  # enable php extensions for a specific app
  nitro ext --app myapp.nitro`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "extensions",
		Short:   "Enables a PHP extension for an app.",
		Example: exampleText,
		Aliases: []string{"ext"},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.RunApply(cmd, args, false, output)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// load the configuration
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// get the app
			appName := flags.AppName
			if appName == "" {
				// get the current working directory
				wd, err := os.Getwd()
				if err != nil {
					return err
				}

				appName, err = appaware.Detect(*cfg, wd)
				if err != nil {
					return err
				}
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			output.Info("modifying", appName)

			// add the label to get the site
			filter.Add("label", containerlabels.Host+"="+appName)

			// find the containers but limited to the site label
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter, All: true})
			if err != nil {
				return err
			}

			// are there any containers??
			if len(containers) == 0 {
				return fmt.Errorf("unable to find an matching site")
			}

			// start the container if it's not running
			if containers[0].State != "running" {
				for _, command := range cmd.Root().Commands() {
					if command.Use == "start" {
						if err := command.RunE(cmd, []string{}); err != nil {
							return err
						}
					}
				}
			}

			// set the hostname of the site based on the container name
			hostname := strings.TrimLeft(containers[0].Names[0], "/")

			extensions := []string{
				"bcmath",
				"bz2",
				"calendar",
				"dba",
				"enchant",
				"exif",
				"gettext",
				"gmp",
				"imap",
				"interbase",
				"ldap",
				"mongo",
				"mysqli",
				"oci8",
				"odbc",
				"pcntl",
				"pdo_dblib",
				"pdo_firebird",
				"pdo_oci",
				"pdo_odbc",
				"pdo_sqlite",
				"recode",
				"shmop",
				"snmp",
				"sockets",
				"sysvmsg",
				"sysvsem",
				"sysvshm",
				"tidy",
				"wddx",
				"xmlrpc",
				"xsl",
				"zend_test",
			}

			// which extensions to add
			selected, err := output.Select(cmd.InOrStdin(), "Which PHP extension would you like to enable for "+hostname+"? ", extensions)
			if err != nil {
				return err
			}

			// get the specific extensions to install
			extension := extensions[selected]

			// set the extension
			if err := cfg.SetPHPExtension(hostname, extension); err != nil {
				return err
			}

			// save the config file
			if err := cfg.Save(); err != nil {
				return fmt.Errorf("unable to save config, %w", err)
			}

			return nil
		},
	}

	return cmd
}
