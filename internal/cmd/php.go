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
		return cmd.Usage()
	},
}

var phpRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "Restart php-fpm",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		runner := nitro.NewMultipassRunner("multipass")
		ip := nitro.IP(machine, runner)
		client := nitrod.NewClient(ip)
		php := config.GetString("php", flagPhpVersion)

		success, err := client.ServicePhpFpm(cmd.Context(), &nitrod.PhpFpmOptions{
			Version: php,
			Action:  "restart",
		})
		if err != nil {
			return err
		}

		fmt.Println(success.Output)

		return nil
	},
}

var phpStartCommand = &cobra.Command{
	Use:   "start",
	Short: "Start php-fpm",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		runner := nitro.NewMultipassRunner("multipass")
		ip := nitro.IP(machine, runner)
		client := nitrod.NewClient(ip)
		php := config.GetString("php", flagPhpVersion)

		success, err := client.ServicePhpFpm(cmd.Context(), &nitrod.PhpFpmOptions{
			Version: php,
			Action:  "start",
		})
		if err != nil {
			return err
		}

		fmt.Println(success.Output)

		return nil
	},
}

var phpStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stop php-fpm",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		runner := nitro.NewMultipassRunner("multipass")
		ip := nitro.IP(machine, runner)
		client := nitrod.NewClient(ip)
		php := config.GetString("php", flagPhpVersion)

		success, err := client.ServicePhpFpm(cmd.Context(), &nitrod.PhpFpmOptions{
			Version: php,
			Action:  "stop",
		})
		if err != nil {
			return err
		}

		fmt.Println(success.Output)

		return nil
	},
}

func init() {
	phpRestartCommand.Flags().StringVar(&flagPhpVersion, "php-version", "", "which PHP version")
	phpStartCommand.Flags().StringVar(&flagPhpVersion, "php-version", "", "which PHP version")
	phpStopCommand.Flags().StringVar(&flagPhpVersion, "php-version", "", "which PHP version")
}
