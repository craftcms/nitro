package hosts

import (
	"strings"

	"github.com/craftcms/nitro/terminal"
	"github.com/spf13/cobra"
)

const exampleText = `  # modify hosts file
  nitro hosts`

// New returns a command used to modify the hosts file to point sites to the nitro
// proxy.
func New(home string, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "hosts",
		Short:   "Modify your hosts file",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			hosts := cmd.Flag("hosts").Value.String()

			// if there are no hosts as flags, use the config file
			if hosts == "" {
			}

			// remove [ and ] from the string
			hosts = strings.Replace(hosts, "[", "", 1)
			hosts = strings.Replace(hosts, "]", "", 1)

			for _, h := range strings.Split(hosts, ",") {
				output.Info(h)
			}

			return nil
		},
	}

	// set flags for the command
	cmd.Flags().StringSliceP("hosts", "z", nil, "hostnames to set")

	return cmd
}
