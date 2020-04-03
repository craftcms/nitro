package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command"
	"github.com/craftcms/nitro/internal/nitro"
)

var siteCommand = &cobra.Command{
	Use:   "site",
	Short: "Perform site commands",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		if err := nitro.Run(
			command.NewMultipassRunner("multipass"),
			nitro.Empty(flagMachineName),
		); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(siteCommand)
}
