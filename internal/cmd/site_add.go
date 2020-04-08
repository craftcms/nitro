package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/validate"
)

var siteAddCommand = &cobra.Command{
	Use:   "add",
	Short: "Add a site to machine",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		machineName := config.GetString("name", flagMachineName)
		php := config.GetString("php", flagPhpVersion)
		localDirectory := args[0]
		domainName := args[1]

		if err := validate.Path(localDirectory); err != nil {
			return err
		}
		if err := validate.Domain(domainName); err != nil {
			return err
		}

		var commands []nitro.Command
		// attach the provided localDirectory to /app/sites/domainName.test
		commands = append(commands, nitro.Mount(machineName, localDirectory, domainName))
		// create localDirectory directory
		commands = append(commands, nitro.CreateNewDirectoryForSite(machineName, domainName))
		// copy the template
		commands = append(commands, nitro.CopyNginxTemplate(machineName, domainName))
		// change template variables
		commands = append(commands, nitro.ChangeVariablesInTemplate(machineName, domainName, flagPublicDir, php)...)
		// make link for nginx localDirectory
		commands = append(commands, nitro.LinkNginxSite(machineName, domainName))
		// reload nginx
		commands = append(commands, nitro.ReloadNginx(machineName))

		if flagDebug {
			for _, command := range commands {
				fmt.Println(command.Type, command.Args)
			}

			return nil
		}

		return nitro.Run(nitro.NewMultipassRunner("multipass"), commands)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println(
			fmt.Sprintf("added site %q to machine %q", args[1], config.GetString("name", flagMachineName)),
		)
	},
}

func init() {
	siteAddCommand.Flags().StringVarP(&flagPhpVersion, "php-version", "p", "", "version of PHP to use")
	siteAddCommand.Flags().StringVarP(&flagPublicDir, "public-dir", "r", "web", "name of the public directory")
}
