package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/client"
	"github.com/craftcms/nitro/internal/config"
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
		c, err := client.NewSystemClient(ip, "50051")
		if err != nil {
			return err
		}
		php := config.GetString("php", flagPhpVersion)

		resp, err := c.PhpFpm(cmd.Context(), &nitrod.PhpFpmServiceRequest{Version: php, Action: nitrod.ServiceAction_RESTART})
		if err != nil {
			return err
		}

		fmt.Println(resp.Message)

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
		c, err := client.NewSystemClient(ip, "50051")
		if err != nil {
			return err
		}
		php := config.GetString("php", flagPhpVersion)

		resp, err := c.PhpFpm(context.WithValue(cmd.Context(), "version", Version), &nitrod.PhpFpmServiceRequest{
			Version: php,
			Action:  nitrod.ServiceAction_START,
		})
		if err != nil {
			return err
		}

		fmt.Println(resp.Message)

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
		c, err := client.NewSystemClient(ip, "50051")
		if err != nil {
			return err
		}
		php := config.GetString("php", flagPhpVersion)

		resp, err := c.PhpFpm(cmd.Context(), &nitrod.PhpFpmServiceRequest{Version: php, Action: nitrod.ServiceAction_STOP})
		if err != nil {
			return err
		}

		fmt.Println(resp.Message)

		return nil
	},
}

func init() {
	phpRestartCommand.Flags().StringVar(&flagPhpVersion, "php-version", "", "which PHP version")
	phpStartCommand.Flags().StringVar(&flagPhpVersion, "php-version", "", "which PHP version")
	phpStopCommand.Flags().StringVar(&flagPhpVersion, "php-version", "", "which PHP version")
}
