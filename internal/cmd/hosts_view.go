package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/txn2/txeh"

	"github.com/craftcms/nitro/config"
)

var hostsShowCommand = &cobra.Command{
	Use:   "view",
	Short: "View your hosts file",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = config.GetString("name", flagMachineName)

		hosts, err := txeh.NewHostsDefault()
		if err != nil {
			return err
		}

		fmt.Println(hosts.RenderHostsFile())

		return nil
	},
}
