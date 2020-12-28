package version

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/craftcms/nitro/protob"
)

// Version is used to set the version of nitro we are using
// and is also used to sync the docker image for the proxy
// container to use to verify the gRPC API is in sync.
var Version = "develop"

var exampleText = `nitro version`

// New is used to show the cli and gRPC API client version
func New(client client.CommonAPIClient, nitrod protob.NitroClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Show version info",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			nitro, err := nitrod.Version(cmd.Context(), &protob.VersionRequest{})
			if err != nil {
				return fmt.Errorf("unable to get version from the gRPC API")
			}

			ver, err := client.ServerVersion(cmd.Context())
			if err != nil {
				return fmt.Errorf("unable to get docker server version, %w", err)
			}

			output.Info(fmt.Sprintf("View the changelog at https://github.com/craftcms/nitro/blob/%s/CHANGELOG.md\n", Version))

			output.Info("Nitro CLI: \t", Version)
			output.Info("Nitro gRPC: \t", nitro.GetVersion())
			output.Info("Docker API: \t", ver.APIVersion, "("+ver.MinAPIVersion+" min)")
			output.Info("Docker CLI: \t", client.ClientVersion())

			// check if the cli and API do not match
			if Version != nitro.GetVersion() {
				output.Info("")
				output.Info("The Nitro CLI and gRPC versions do not match")
				output.Info("You might need to run `nitro update`")
			}

			return nil
		},
	}

	return cmd
}
