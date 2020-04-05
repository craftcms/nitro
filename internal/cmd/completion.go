package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

. <(nitro completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(nitro completion)
`,
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nitroCommand.GenBashCompletion(os.Stdout)
	},
}
