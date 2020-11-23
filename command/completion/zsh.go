package completion

import (
	"fmt"
	"os"

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
$ nitro complete zsh > "${fpath[1]}/_nitro"

# You will need to start a new shell for this setup to take effect.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.GenZshCompletion(os.Stdout)
		if err != nil {
			fmt.Println(err)
		}
	},
}
