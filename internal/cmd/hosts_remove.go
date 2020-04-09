package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/txn2/txeh"

	"github.com/craftcms/nitro/config"
)

var hostsRemoveCommand = &cobra.Command{
	Use:   "remove",
	Short: "Remove entries from /etc/hosts file",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = config.GetString("name", flagMachineName)

		uid := os.Geteuid()
		if uid != 0 {
			return errors.New("you do not appear to be running this command as root, so we cannot modify the /etc/hosts file")
		}

		// get all of the sites from the config file
		if !viper.IsSet("sites") {
			return errors.New("unable to read sites from " + viper.ConfigFileUsed())
		}

		var sites []config.Site
		if err := viper.UnmarshalKey("sites", &sites); err != nil {
			return err
		}

		hosts, err := txeh.NewHostsDefault()
		if err != nil {
			return err
		}

		var domains []string
		for _, site := range sites {
			domains = append(domains, site.Domain)
		}

		hosts.RemoveHosts(domains)

		return hosts.Save()
	},
}
