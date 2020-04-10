package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
)

var siteRemoveCommand = &cobra.Command{
	Use:   "remove",
	Short: "Remove a site from machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = config.GetString("name", flagMachineName)

		return errors.New("not implemented yet")
	},
}
