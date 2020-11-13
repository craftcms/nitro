package trust

import (
	"fmt"

	"github.com/craftcms/nitro/pkg/client"
	"github.com/spf13/cobra"
)

// TrustCommand is used to retreive the certificates from the nitro proxy and install on the
// host machine. The CA is used to sign certificates for websites and adding the certificate
// to the system allows TLS connections to be considered valid and trusted from the container.
var TrustCommand = &cobra.Command{
	Use:   "trust",
	Short: "Trust certificates for environment",
	RunE:  start,
	Example: `  # get the certificates for an environment
  nitro trust`,
}

func start(cmd *cobra.Command, args []string) error {
	env := cmd.Flag("environment").Value.String()

	// create the new client
	nitro, err := client.NewClient()
	if err != nil {
		return fmt.Errorf("unable to create a client for docker, %w", err)
	}

	return nitro.Trust(cmd.Context(), env, args)
}
