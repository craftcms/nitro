package portcheck

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/portavail"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # check if a port is in use
  nitro portcheck 8080`

// NewCommand returns the command to enable common nitro services. These services are provided as containers
// and do not require a user to configure the ports/volumes or images.
func NewCommand(output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "portcheck",
		Short:   "Check a local port",
		Args:    cobra.MinimumNArgs(1),
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := portavail.Check(args[0]); err != nil {
				output.Info("Port", args[0], "is already in use...")

				return nil
			}

			output.Info("Port", args[0], "is available!")

			return nil
		},
	}

	return cmd
}
