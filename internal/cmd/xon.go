package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/client"
	"github.com/craftcms/nitro/internal/config"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/nitrod"
)

func init() {
	xonCommand.Flags().StringVar(&flagPhpVersion, "php-version", "", "which PHP version")
}

var xonCommand = &cobra.Command{
	Use:   "xon",
	Short: "Enable Xdebug",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		runner := nitro.NewMultipassRunner("multipass")
		ip := nitro.IP(machine, runner)
		c, err := client.NewClient(ip, "50051")
		if err != nil {
			return err
		}
		php := config.GetString("php", flagPhpVersion)

		resp, err := c.EnableXdebug(cmd.Context(), &nitrod.EnableXdebugRequest{Version: php})
		if err != nil {
			return err
		}

		if !flagSilent {
			fmt.Println(resp.Message)
		}

		return nil
	},
}

func init() {
	xonCommand.Flags().BoolVar(&flagSilent, "silent", false, "Run command with no output")
}
