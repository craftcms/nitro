package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/nitrod"
	"github.com/craftcms/nitro/internal/client"
	"github.com/craftcms/nitro/internal/nitro"
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
		c := client.NewSystemClient(ip, "50051")

		resp, err := c.Nginx(cmd.Context(), &nitrod.NginxServiceRequest{Action: nitrod.ServiceAction_RESTART})
		if err != nil {
			return err
		}

		fmt.Println(resp.Message)

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
		c := client.NewSystemClient(ip, "50051")

		resp, err := c.Nginx(cmd.Context(), &nitrod.NginxServiceRequest{Action: nitrod.ServiceAction_START})
		if err != nil {
			return err
		}

		fmt.Println(resp.Message)

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
		c := client.NewSystemClient(ip, "50051")

		resp, err := c.Nginx(cmd.Context(), &nitrod.NginxServiceRequest{Action: nitrod.ServiceAction_STOP})
		if err != nil {
			return err
		}

		fmt.Println(resp.Message)

		return nil
	},
}
