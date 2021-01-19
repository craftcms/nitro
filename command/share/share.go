package share

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	// ErrMissingNgrok is returned when nirto is unable to locate the ngrok binary
	ErrMissingNgrok = fmt.Errorf("unable to locate ngrok")
)

const exampleText = `  # share a local site with ngrok
  nitro share`

// NewCommand is used to destroy all resources for an environment. It will prompt for
// user verification and defaults to no. Part of the destroy process is to
// perform a backup for all databases in each container database.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "share",
		Short:   "Share a local site",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// create a filter for the environment
			filter := filters.NewArgs()
			filter.Add("label", labels.Nitro)

			// get all of the sites
			var sites, found []string
			for _, s := range cfg.Sites {
				p, _ := s.GetAbsPath(home)

				// check if the path matches a sites path, then we are in a known site
				if strings.Contains(wd, p) {
					found = append(found, s.Hostname)
				}

				// add the site to the list in case we cannot find the directory
				sites = append(sites, s.Hostname)
			}

			// if there are found sites we want to show or connect to the first one, otherwise prompt for
			// which site to connect to.
			switch len(found) {
			case 0:
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", sites)
				if err != nil {
					return err
				}

				// add the label to get the site
				filter.Add("label", labels.Host+"="+sites[selected])
			case 1:
				// add the label to get the site
				filter.Add("label", labels.Host+"="+found[0])
			default:
				// prompt for the site to ssh into
				selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", found)
				if err != nil {
					return err
				}

				// add the label to get the site
				filter.Add("label", labels.Host+"="+found[selected])
			}

			// find the containers but limited to the site label
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter})
			if err != nil {
				return err
			}

			// are there any containers??
			if len(containers) == 0 {
				return fmt.Errorf("unable to find an matching site")
			}

			// start the container if its not running
			if containers[0].State != "running" {
				if err := docker.ContainerStart(ctx, containers[0].ID, types.ContainerStartOptions{}); err != nil {
					return err
				}
			}

			// find ngrok
			ngrok, err := exec.LookPath("ngrok")
			if err != nil {
				return err
			}

			if ngrok == "" {
				return ErrMissingNgrok
			}

			hostname := strings.TrimLeft(containers[0].Names[0], "/")

			c := exec.Command(ngrok, "http", "-host-header="+hostname, "80")

			c.Stderr = cmd.ErrOrStderr()
			c.Stdout = cmd.OutOrStdout()

			return c.Run()
		},
	}

	// add flags to the command
	cmd.Flags().Bool("clean", false, "remove configuration file")

	return cmd
}
