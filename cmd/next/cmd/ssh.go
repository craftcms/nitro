package cmd

import (
	"github.com/spf13/cobra"
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		//if err := nitro.Run(
		//	command.NewMultipassRunner("multipass"),
		//	nitro.SSH(flagMachineName),
		//); err != nil {
		//	log.Fatal(err)
		//}
	},
}
