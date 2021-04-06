package completion

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const exampleText = `To load completions:

Bash:

$ source <(nitro completion bash)

# To load completions for each session, execute once:
Linux:
  $ nitro completion bash | sudo tee -a /etc/bash_completion.d/nitro >/dev/null
MacOS:
  $ nitro completion bash > /usr/local/etc/bash_completion.d/nitro

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ nitro completion zsh > "${fpath[1]}/_nitro"

# You will need to start a new shell for this setup to take effect.
`

// NewCommand returns the command used for generating completion shells
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "completion",
		Short:     "Enables shell completion.",
		ValidArgs: []string{"bash", "zsh"},
		Example:   exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// print the help if not defined
			if len(args) == 0 {
				return cmd.Help()
			}

			switch args[0] {
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			}

			return fmt.Errorf("unknown shell requested")
		},
	}

	return cmd
}
