package completion

import (
	"fmt"

	"github.com/spf13/cobra"
)

// zshCompletionCommand represents the completion command for zsh
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
		// Gross, but fixes a ZSH completion bug where Cobra generates the completion script using the command name.
		cmd.Use = "nitro"
		err := cmd.GenZshCompletion(cmd.OutOrStdout())
		cmd.Use = "zsh"
		if err != nil {
			fmt.Println(err)
		}
	},
}
