package extensions

import (
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	// ErrUnknownExtension is used when an unknown service is requested
	ErrUnknownExtension = fmt.Errorf("unknown extension requested")
)

const exampleText = `  # enable PHP extensions for a site
  nitro extensions`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "extensions",
		Short:   "Enables a PHP extension for a site.",
		Example: exampleText,
		Aliases: []string{"ext"},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.RunApply(cmd, args, false, output)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// load the configuration
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			// get a context aware list of sites
			sites := cfg.ListOfSitesByDirectory(home, wd)

			// create the options for the sites
			var options []string
			for _, s := range sites {
				options = append(options, s.Hostname)
			}

			// if there are found sites we want to show or connect to the first one, otherwise prompt for
			// which site to connect to.
			switch len(sites) {
			case 0:
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
				if err != nil {
					return err
				}

				// add the label to get the site
				filter.Add("label", containerlabels.Host+"="+sites[selected].Hostname)
			case 1:
				output.Info("modifiying", sites[0].Hostname)

				// add the label to get the site
				filter.Add("label", containerlabels.Host+"="+sites[0].Hostname)
			default:
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
				if err != nil {
					return err
				}

				// add the label to get the site
				filter.Add("label", containerlabels.Host+"="+sites[selected].Hostname)
			}

			// find the containers but limited to the site label
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter, All: true})
			if err != nil {
				return err
			}

			// are there any containers??
			if len(containers) == 0 {
				return fmt.Errorf("unable to find an matching site")
			}

			// start the container if its not running
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
