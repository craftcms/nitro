package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/validate"
)

var (
	flagPublicDir string

	siteCommand = &cobra.Command{
		Use:   "site",
		Short: "Perform site commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	siteAddCommand = &cobra.Command{
		Use:   "add",
		Short: "Add a site to machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := config.GetString("machine", flagMachineName)
			php := config.GetString("php", flagPhpVersion)
			path := args[0]
			domain := args[1]

			if err := validate.Path(path); err != nil {
				return err
			}

			if err := validate.Domain(domain); err != nil {
				return err
			}

			var commands []nitro.Command

			// attach the provided path to /app/sites/domain.test
			commands = append(commands, nitro.Mount(name, path, domain))

			// run the nginx add-site script
			commands = append(commands, nitro.AddSiteScript(name, domain, php, flagPublicDir))

			// todo edit the hosts file

			if flagDebug {
				for _, command := range commands {
					fmt.Println(command.Type, command.Args)
				}

				return nil
			}

			if err := nitro.Run(nitro.NewMultipassRunner("multipass"), commands); err != nil {
				return err
			}

			return nil
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			fmt.Println(
				fmt.Sprintf("added site %q to machine %q", args[1], config.GetString("machine", flagMachineName)),
			)
		},
	}
)

func init() {
	siteAddCommand.Flags().StringVarP(&flagPhpVersion, "php-version", "p", "", "version of PHP to use")
	siteAddCommand.Flags().StringVarP(&flagPublicDir, "public-dir", "r", "web", "name of the public directory (defaults to web)")
}
