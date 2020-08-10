package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/config"
	"github.com/craftcms/nitro/internal/nitro"
)

const infoTemplate = `Nitro installed, ready for something incredible at %s.

Add a project with "nitro add <directory>". 

Server Information
-------------------------
IP address: %s
PHP version: %s

Need help setting up Xdebug?
https://craftcms.com/docs/nitro/xdebug.html

Need help using Nitro? 
https://craftcms.com/docs/nitro`

var infoCommand = &cobra.Command{
	Use:   "info",
	Short: "Show machine info",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		runner := nitro.NewMultipassRunner("multipass")
		ip := nitro.IP(machine, runner)
		php := config.GetString("php", flagPhpVersion)

		// check if the machine is running, if not start it
		if ip == "" {
			fmt.Println("The " + machine + " machine is not running...")
			if err := startCommand.RunE(cmd, args); err != nil {
				return err
			}

			// get the IP again
			ip = nitro.IP(machine, runner)
		}

		fmt.Println(fmt.Sprintf(infoTemplate, ip, ip, php))
		fmt.Println("")

		return nil
	},
}
