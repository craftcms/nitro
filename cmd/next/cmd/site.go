package cmd

import (
	"github.com/spf13/cobra"
)

var siteCommand = &cobra.Command{
	Use:   "site",
	Short: "Perform site commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		//if len(args) == 0 {
		//	_ = cmd.Help()
		//	return
		//}

		return nil

		//if err := nitro.Run(
		//	command.NewMultipassRunner("multipass"),
		//	nitro.Empty(flagMachineName),
		//); err != nil {
		//	log.Fatal(err)
		//}
	},
}
