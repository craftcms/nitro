package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/nitrod"
)

var nginxCommand = &cobra.Command{
	Use:   "nginx",
	Short: "Perform nginx actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

var nginxRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "Restart nginx",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		runner := nitro.NewMultipassRunner("multipass")
		ip := nitro.IP(machine, runner)
		client := nitrod.NewClient(ip)

		success, err := client.ServiceNginx(cmd.Context(), &nitrod.NginxServiceOptions{
			Action: "restart",
		})
		if err != nil {
			return err
		}

		fmt.Println(success.Output)

		return nil
	},
}

var nginxStartCommand = &cobra.Command{
	Use:   "start",
	Short: "Start nginx",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		runner := nitro.NewMultipassRunner("multipass")
		ip := nitro.IP(machine, runner)
		client := nitrod.NewClient(ip)

		success, err := client.ServiceNginx(cmd.Context(), &nitrod.NginxServiceOptions{
			Action: "start",
		})
		if err != nil {
			return err
		}

		fmt.Println(success.Output)

		return nil
	},
}

var nginxStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stop nginx",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		runner := nitro.NewMultipassRunner("multipass")
		ip := nitro.IP(machine, runner)
		client := nitrod.NewClient(ip)

		success, err := client.ServiceNginx(cmd.Context(), &nitrod.NginxServiceOptions{
			Action: "stop",
		})
		if err != nil {
			return err
		}

		fmt.Println(success.Output)

		return nil
	},
}
