package hosts

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/craftcms/nitro/hostedit"
	"github.com/craftcms/nitro/terminal"
	"github.com/spf13/cobra"
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
			hosts := cmd.Flag("hostnames").Value.String()
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

			// set the file based on the OS
			defaultFile := "/etc/hosts"
			if runtime.GOOS == "windows" {
				defaultFile = `C:\Windows\System32\Drivers\etc\hosts`
			}

			// is this a remove or add
			var updatedContent string
			switch cmd.Flag("remove").Value.String() {
			case "true":
				if !preview {
					output.Info("Removing sites from hosts file...")
				}
				return fmt.Errorf("remove is not yet implemented")
			default:
				if !preview {
					output.Info("Adding sites to hosts file...")
				}

				// add the hosts
				updated, err := hostedit.Update(defaultFile, "127.0.0.1", hostnames...)
				if err != nil {
					return err
				}
				updatedContent = updated
			}

			// if we are previewing, show the hosts file without saving
			if preview {
				output.Info("Previewing changes to hosts file...\n")

				output.Info(updatedContent)

				return nil
			}

			// check if we are the root user
			uid := os.Geteuid()
			if (uid != 0) && (uid != -1) {
				return fmt.Errorf("you do not appear to be running this command as root, so we cannot modify your hosts file")
			}

			output.Pending("modifying hosts file")

			// save the file
			if err := ioutil.WriteFile(defaultFile, []byte(updatedContent), 0644); err != nil {
				return err
			}

			output.Done()

			return nil
		},
	}

	// set flags for the command
	cmd.Flags().StringSliceP("hostnames", "z", nil, "list of hostnames to set")
	cmd.MarkFlagRequired("hostnames")
	cmd.Flags().BoolP("remove", "r", false, "remove hosts from file")
	cmd.Flags().BoolP("preview", "p", false, "preview the hosts file change")

	return cmd
}
