package keys

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/keys"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/spf13/cobra"
)

var (
	// when a site is selected, we use a global variable to keep the code clean
	site *config.Site
)

const exampleText = `  # keys command
  nitro keys`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys",
		Short:   "Adds SSH keys to a site container.",
		Example: exampleText,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			cfg, err := config.Load(home)
			if err != nil {
				return nil, cobra.ShellCompDirectiveDefault
			}

			var options []string
			for _, s := range cfg.Sites {
				options = append(options, s.Hostname)
			}

			return options, cobra.ShellCompDirectiveNoFileComp
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// is there a site as the first arg?
			if len(args) > 0 {
				site, err = cfg.FindSiteByHostName(args[0])
				if err != nil {
					return err
				}

				output.Info("Preparing key import to", site.Hostname)

				return nil
			}

			// get a context aware list of sites
			sites := cfg.ListOfSitesByDirectory(home, wd)

			// create the options for the sites
			var options []string
			for _, s := range sites {
				options = append(options, s.Hostname)
			}

			// prompt for the site to ssh into
			selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
			if err != nil {
				return err
			}

			site = &sites[selected]

			output.Info("Preparing key import to", site.Hostname)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			path := filepath.Join(home, ".ssh")
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return errors.New("unable to find directory " + path)
			}

			// create a filter
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)
			filter.Add("label", containerlabels.Host+"="+site.Hostname)

			// find the containers but limited to the site label
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{Filters: filter, All: true})
			if err != nil {
				return err
			}

			// make sure there are
			if len(containers) == 0 {
				return fmt.Errorf("no containers found")
			}

			// set the container
			container := containers[0]

			// start the container if it's not running
			if container.State != "running" {
				output.Pending("starting container for", site.Hostname)

				// start the container
				if err := docker.ContainerStart(cmd.Context(), container.ID, types.ContainerStartOptions{}); err != nil {
					return err
				}

				output.Done()
			}

			// find all keys
			found, err := keys.Find(path)
			if err != nil {
				return err
			}

			var options []string
			for k := range found {
				options = append(options, k)
			}

			// prompt the user for their selected key
			selected, err := output.Select(os.Stdin, "Which key should we add to the site?", options)
			if err != nil {
				return err
			}

			// key use a map for the keys, so we have this loop
			var count int
			var key string
			for k := range found {
				if count == selected {
					key = k
					break
				}

				count++
			}

			// verify the key (using the docker stat path API) does not already exist in /home/nitro/.ssh/<key>
			stat, err := docker.ContainerStatPath(cmd.Context(), container.ID, fmt.Sprintf("/home/nitro/.ssh/%s", found[key]))
			if err != nil {
				// the docker sdk does not return an error, so we have to check the error output
				if strings.Contains(err.Error(), "Error: No such container:path") {
					// do nothing
				} else {
					return err
				}
			}

			// check if the file exists
			if stat.Name != "" {
				// prompt the user to confirm overwriting the file
				confirm, err := output.Confirm("The key appears to exist in the site, should we overwrite", true, "?")
				if err != nil {
					return err
				}

				if !confirm {
					output.Info("Skipping")
					return nil
				}
			}

			content, err := ioutil.ReadFile(filepath.Join(path, found[key]))
			if err != nil {
				return err
			}

			// create the temp file
			tr, err := archive.Generate(string(content))
			if err != nil {
				return err
			}

			// copy the file into the container
			if err := docker.CopyToContainer(cmd.Context(), container.ID, "/home/nitro/.ssh/"+found[key], tr, types.CopyToContainerOptions{AllowOverwriteDirWithFile: false}); err != nil {
				fmt.Println("OOPS")
				return err
			}

			return nil
		},
	}

	return cmd
}
