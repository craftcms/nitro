package hosts

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/craftcms/nitro/terminal"
	"github.com/spf13/cobra"
	"github.com/txn2/txeh"
)

const exampleText = `  # modify hosts file to match sites and aliases
  nitro hosts

  # remove sites and aliases from hosts file
  nitro hosts --remove`

// New returns a command used to modify the hosts file to point sites to the nitro
// proxy.
func New(home string, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "hosts",
		Short:   "Modify your hosts file",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			hosts := cmd.Flag("hosts").Value.String()
			preview, err := strconv.ParseBool(cmd.Flag("preview").Value.String())
			if err != nil {
				// if there is an error set preview to false
				preview = false
			}

			var hostnames []string

			// remove [ and ] from the string
			hosts = strings.Replace(hosts, "[", "", 1)
			hosts = strings.Replace(hosts, "]", "", 1)

			for _, h := range strings.Split(hosts, ",") {
				hostnames = append(hostnames, h)
			}

			// create the host editor, should be a dependency to the function
			hostedit, err := txeh.NewHostsDefault()
			if err != nil {
				return err
			}

			// is this a remove or add
			switch cmd.Flag("remove").Value.String() {
			case "true":
				if !preview {
					output.Info("Removing sites from hosts file...")
				}

				// remove the hosts from the file
				hostedit.RemoveHosts(hostnames)
			default:
				if !preview {
					output.Info("Adding sites to hosts file...")
				}

				hostedit.AddHosts("127.0.0.1", hostnames)
			}

			// if we are previewing, show the hosts file without saving
			if preview {
				output.Info("Previewing changes to hostfile...\n")

				output.Info(hostedit.RenderHostsFile())

				return nil
			}

			// check if we are the root user
			uid := os.Geteuid()
			if (uid != 0) && (uid != -1) {
				return fmt.Errorf("you do not appear to be running this command as root, so we cannot modify your hosts file")
			}

			output.Pending("modifying hosts file")

			// try to save the hosts file
			if err := hostedit.Save(); err != nil {
				output.Warning()
				return err
			}

			output.Done()

			return nil
		},
	}

	// set flags for the command
	cmd.Flags().StringSliceP("hosts", "z", nil, "list of hostnames to set")
	cmd.MarkFlagRequired("hosts")
	cmd.Flags().BoolP("remove", "r", false, "remove hosts from file")
	cmd.Flags().BoolP("preview", "p", false, "preview the hosts file change")

	return cmd
}
