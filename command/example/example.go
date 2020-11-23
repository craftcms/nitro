package example

import (
	"fmt"

	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var (
	// ErrExample is used when we want to share an error
	ErrExample = fmt.Errorf("some example error")
)

const exampleText = `  # example command
  nitro example`

// New is used for scaffolding new commands
func New(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "example",
		Short:   "Example command",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			return ErrExample
		},
	}

	// set flags for the command
	cmd.Flags().StringP("example", "e", "example", "an example flag")

	return cmd
}
