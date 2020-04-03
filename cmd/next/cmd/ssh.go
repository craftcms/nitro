package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command"
	"github.com/craftcms/nitro/internal/nitro"
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		if err := nitro.Run(
			command.NewMultipassRunner("multipass"),
			nitro.SSH(flagMachineName),
		); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
}
