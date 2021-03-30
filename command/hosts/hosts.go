package hosts

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/hostedit"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # modify hosts file to match sites and aliases
  nitro hosts`

// New returns a command used to modify the hosts file to point sites to the nitro proxy.
func NewCommand(home string, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "hosts",
		Short:   "Modifies hosts file.",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			hosts := cmd.Flag("hostnames").Value.String()
			var preview bool
			previewFlag := cmd.Flag("preview").Value.String()
			if previewFlag == "true" {
				preview = true
			}

			// remove [ and ] from the string
			hosts = strings.Replace(hosts, "[", "", 1)
			hosts = strings.Replace(hosts, "]", "", 1)

			var hostnames []string
			hostnames = append(hostnames, strings.Split(hosts, ",")...)

			// set the file based on the OS
			defaultFile := "/etc/hosts"
			if runtime.GOOS == "windows" {
				defaultFile = `C:\Windows\System32\Drivers\etc\hosts`
			}

			// add the hosts
			updated, err := hostedit.Update(defaultFile, "127.0.0.1", hostnames...)
			if err != nil {
				return err
			}

			// if we are previewing, show the hosts file without saving
			if preview {
				output.Info("Previewing changes to hosts file…\n")

				output.Info(updated)

				return nil
			} else {
				output.Info("Adding sites to hosts file…")
			}

			// check if we are the root user
			uid := os.Geteuid()
			if (uid != 0) && (uid != -1) {
				return fmt.Errorf("you do not appear to be running this command as root, so we cannot modify your hosts file")
			}

			output.Pending("modifying hosts file")

			// save the file
			if err := ioutil.WriteFile(defaultFile, []byte(updated), 0644); err != nil {
				return err
			}

			output.Done()

			return nil
		},
	}

	// set flags for the command
	cmd.Flags().StringSlice("hostnames", nil, "list of hostnames to set")
	cmd.MarkFlagRequired("hostnames")
	cmd.Flags().Bool("preview", false, "preview hosts file change")

	cmd.AddCommand(removeCommand(home, output))

	return cmd
}
