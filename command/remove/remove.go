package remove

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # remove a site from the config
  nitro remove`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Short:   "Remove a site",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// get all of the sites
			var sites []string
			for _, s := range cfg.Sites {
				// add the site to the list
				sites = append(sites, s.Hostname)
			}

			// prompt for the site to remove
			selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", sites)
			if err != nil {
				return err
			}

			if err := cfg.RemoveSite(sites[selected]); err != nil {
				return err
			}

			// ask if the apply command should run
			fmt.Print("Apply changes now [Y/n]? ")

			s := bufio.NewScanner(os.Stdin)
			s.Split(bufio.ScanLines)

			var confirm bool
			for s.Scan() {
				txt := strings.TrimSpace(s.Text())

				switch txt {
				// if its no
				case "n", "N", "no", "No", "NO":
					confirm = false
				default:
					confirm = true
				}

				break
			}
			if err := s.Err(); err != nil {
				return err
			}

			// we are skipping the apply step
			if !confirm {
				return nil
			}

			// check if there is no parent command
			if cmd.Parent() == nil {
				return nil
			}

			// get the apply command and run it
			for _, c := range cmd.Parent().Commands() {
				if c.Use == "apply" {
					return c.RunE(c, args)
				}
			}

			return nil
		},
	}

	// cmd.Flags().Bool("proxy", false, "connect to the proxy container")

	return cmd
}
