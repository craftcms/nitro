package complete

import "github.com/spf13/cobra"

// CompleteCommand is the command for shell completion
var CompleteCommand = &cobra.Command{
	Use:   "complete",
	Short: "Enable shell completion",
	RunE:  completeMain,
	Long: `To load completions:

Bash:

$ source <(nitro completion bash)

# To load completions for each session, execute once:
Linux:
  $ nitro complete bash > /etc/bash_completion.d/nitro
MacOS:
  $ nitro complete bash > /usr/local/etc/bash_completion.d/nitro

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ nitro complete zsh > "${fpath[1]}/_nitro"

# You will need to start a new shell for this setup to take effect.
`,
}

func completeMain(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func init() {
	CompleteCommand.AddCommand(bashCompletionCommand, zshCompletionCommand)
}
