package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/nitrod"
)

var phpCommand = &cobra.Command{
	Use:   "php",
	Short: "Perform PHP actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		runner := nitro.NewMultipassRunner("multipass")
		ip := nitro.IP(machine, runner)
		client := nitrod.NewClient(ip)
		php := config.GetString("php", flagPhpVersion)
		action := ""

		if flagRestart {
			action = "restart"
		}

		if flagStop {
			action = "stop"
		}

		if flagStart {
			action = "start"
		}

		switch action {
		case "":
			return cmd.Usage()
		default:
			success, err := client.ServicePhpFpm(cmd.Context(), &nitrod.PhpFpmOptions{
				Version: php,
				Action:  action,
			})
			if err != nil {
				return err
			}

			fmt.Println(success.Output)

			return nil
		}
	},
}

func init() {
	phpCommand.Flags().StringVar(&flagPhpVersion, "php-version", "", "which PHP version")
	phpCommand.Flags().BoolVar(&flagRestart, "restart", false, "restart PHP-FPM")
	phpCommand.Flags().BoolVar(&flagStop, "stop", false, "stop PHP-FPM")
	phpCommand.Flags().BoolVar(&flagStart, "start", false, "start PHP-FPM")
}
