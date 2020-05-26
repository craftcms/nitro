package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/txn2/txeh"

	"github.com/craftcms/nitro/internal/hosts"
)

var hostsRemoveCommand = &cobra.Command{
	Use:    "remove",
	Short:  "Remove site from your hosts file",
	Hidden: true,
	Args:   cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !flagDebug {
			uid := os.Geteuid()
			if uid != 0 {
				return errors.New("you do not appear to be running this command as root, so we cannot modify your hosts file")
			}
		}

		he, err := txeh.NewHostsDefault()
		if err != nil {
			return err
		}

		var domains []string
		for _, site := range args {
			domains = append(domains, site)
		}

		if flagDebug {
			for _, domain := range domains {
				fmt.Println("Removing", domain, "from hosts file.")
			}

			return nil
		}

		return hosts.Remove(he, domains)
	},
}
