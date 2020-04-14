package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var ipCommand = &cobra.Command{
	Use:    "ip",
	Short:  "Show machine IP",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)
		r := nitro.NewMultipassRunner("multipass")

		ip := nitro.IP(name, r)
		if ip == "" {
			return errors.New("could not get the IP of " + name)
		}

		fmt.Println(ip)

		return nil
	},
}
