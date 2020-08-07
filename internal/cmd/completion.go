package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate completion scripts",
	Long: `To load completions:

Bash:

$ source <(nitro completion bash)

# To load completions for each session, execute once:
Linux:
  $ nitro completion bash > /etc/bash_completion.d/nitro
MacOS:
  $ nitro completion bash > /usr/local/etc/bash_completion.d/nitro

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ nitro completion zsh > "${fpath[1]}/_nitro"

# You will need to start a new shell for this setup to take effect.
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
	Short: "Generates zsh completion scripts",
	Long: `To load zsh completion:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ nitro completion zsh > "${fpath[1]}/_nitro"

# You will need to start a new shell for this setup to take effect.
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
		err := rootCmd.GenBashCompletion(os.Stdout)
		if err != nil {
			fmt.Println(err)
		}
	},
}
