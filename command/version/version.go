package version

import (
	"fmt"

	"github.com/craftcms/nitro/protob"
	"github.com/craftcms/nitro/terminal"
	"github.com/spf13/cobra"
)

// Version is used to set the version of nitro we are using
// and is also used to sync the docker image for the proxy
// container to use to verify the gRPC API is in sync.
var Version = "dev"

const exampleText = `  # show the cli and api version
  nitro version`

// New is used to show the cli and gRPC API client version
func New(nitro protob.NitroClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Show version info",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := nitro.Version(cmd.Context(), &protob.VersionRequest{})
			if err != nil {
				return fmt.Errorf("unable to get the version from the gRPC API")
			}

			output.Info("CLI Version:\t", Version)
			output.Info("gRPC Version:\t", resp.GetVersion())

			return nil
		},
	}

	return cmd
}
