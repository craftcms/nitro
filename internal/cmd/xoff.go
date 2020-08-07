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
	xoffCommand.Flags().StringVar(&flagPhpVersion, "php-version", "", "which PHP version")
}

var xoffCommand = &cobra.Command{
	Use:   "xoff",
	Short: "Disable Xdebug",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		runner := nitro.NewMultipassRunner("multipass")
		ip := nitro.IP(machine, runner)
		c, err := client.NewClient(ip, "50051")
		if err != nil {
			return err
		}
		php := config.GetString("php", flagPhpVersion)

		resp, err := c.DisableXdebug(cmd.Context(), &nitrod.DisableXdebugRequest{Version: php})
		if err != nil {
			return err
		}

		fmt.Println(resp.Message)

		return nil
	},
}
