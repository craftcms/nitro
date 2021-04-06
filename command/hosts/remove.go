package hosts

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/craftcms/nitro/pkg/hostedit"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/spf13/cobra"
)

// removeCommand returns a command used to remove entries from the hosts file
func removeCommand(home string, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Removes Nitro entries from hosts file.",
		Example: `  # remove nitro entries from your hosts file
  nitro hosts remove`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var preview bool
			previewFlag := cmd.Flag("preview").Value.String()
			if previewFlag == "true" {
				preview = true
			}

			// set the file based on the OS
			defaultFile := "/etc/hosts"
			if runtime.GOOS == "windows" {
				defaultFile = `C:\Windows\System32\Drivers\etc\hosts`
			}

			// add the hosts
			updated, err := hostedit.Remove(defaultFile)
			if errors.Is(err, hostedit.ErrNotNitroEntries) {
				output.Info("There are no entries to remove from the hosts file...")

				return nil
			}
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
	cmd.Flags().Bool("preview", false, "preview hosts file change")

	return cmd
}
