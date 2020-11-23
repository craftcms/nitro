package completion

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// bashCompletionCommand represents the completion command
var bashCompletionCommand = &cobra.Command{
	Use:   "bash",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

Bash:

$ source <(nitro completion bash)

# To load completions for each session, execute once:
Linux:
  $ nitro completion bash > /etc/bash_completion.d/nitro
MacOS:
  $ nitro completion bash > /usr/local/etc/bash_completion.d/nitro

# You will need to start a new shell for this setup to take effect.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Parent().Parent().GenBashCompletion(os.Stdout)
		if err != nil {
			fmt.Println(err)
		}
	},
}
