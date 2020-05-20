package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate completion scripts",
	Hidden: true,
	Long: `To load completion for bash run

. <(nitro completion bash)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(nitro completion bash)

If you are using ZSH, run

. <(nitro completion zsh)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.zshrc
. <(nitro completion zsh)
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	completionCmd.AddCommand(bashCompletionCommand, zshCompletionCommand)
}

// completionCmd represents the completion command
var zshCompletionCommand = &cobra.Command{
	Use:   "zsh",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

. <(nitro completion zsh)

To configure your zsh shell to load completions for each session add to your zshrc

# ~/.zshrc
. <(nitro completion zsh)
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := rootCmd.GenZshCompletion(os.Stdout)
		if err != nil {
			fmt.Println(err)
		}
	},
}

// completionCmd represents the completion command
var bashCompletionCommand = &cobra.Command{
	Use:   "bash",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

. <(nitro completion bash)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(nitro completion bash)
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := rootCmd.GenBashCompletion(os.Stdout)
		if err != nil {
			fmt.Println(err)
		}
	},
}
