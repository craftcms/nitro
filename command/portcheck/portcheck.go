package portcheck

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/portavail"
	"github.com/craftcms/nitro/pkg/terminal"
)

// Hostname is the default hostname to use for portchecks
var Hostname = "localhost"

const exampleText = `  # check if a port is in use
  nitro portcheck 8080

  # check a port on a specific hostname
  nitro portcheck 8080 --hostname 192.168.7.241`

// NewCommand returns the command to enable common nitro services. These services are provided as containers
// and do not require a user to configure the ports/volumes or images.
func NewCommand(output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "portcheck",
		Short:   "Check a local port",
		Args:    cobra.MinimumNArgs(1),
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// get the port from the args
			port := args[0]

			// check if the port is in use
			if err := portavail.Check(Hostname, port); err != nil {
				output.Info("Port", port, "is already in use...")

				return nil
			}

			output.Info("Port", port, "is available!")

			return nil
		},
	}

	cmd.Flags().StringVar(&Hostname, "hostname", "localhost", "The hostname to use when checking the port")

	return cmd
}
