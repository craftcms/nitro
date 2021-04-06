package permissions

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/terminal"
)

type file struct {
	path      string
	directory bool
}

const exampleText = `  # fix file permissions for craft
  nitro permissions`

var files = []file{
	{
		path:      "/app/.env",
		directory: false,
	},
	{
		path:      "/app/composer.json",
		directory: false,
	},
	{
		path:      "/app/composer.lock",
		directory: false,
	},
	{
		path:      "/app/config/licence.key",
		directory: false,
	},
	{
		path:      "/app/config",
		directory: true,
	},
	{
		path:      "/app/config/project",
		directory: true,
	},
	{
		path:      "/app/storage",
		directory: true,
	},
	{
		path:      "/app/vendor",
		directory: true,
	},
	{
		path:      "/app/web/cpresources",
		directory: true,
	},
}

// NewCommand is the permissions command that sets the proper permissions for files and directories in a container for a site.
// It is primarily used on Linux as Docker Desktop handles permissions for mounts on macOS and Windows.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "permissions",
		Short:   "Fixes Craft permissions.",
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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

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

			// get a context aware list of sites
			sites := cfg.ListOfSitesByDirectory(home, wd)

			// create the options for the sites
			var options []string
			for _, s := range sites {
				options = append(options, s.Hostname)
			}

			var siteArg string
			if len(args) > 0 {
				siteArg = strings.TrimSpace(args[0])
			}

			var site *config.Site
			switch siteArg == "" {
			case true:
				switch len(sites) {
				case 0:
					selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
					if err != nil {
						return err
					}

					site = &sites[selected]
				case 1:
					site = &sites[0]
				default:
					selected, err := output.Select(cmd.InOrStdin(), "Select a site: ", options)
					if err != nil {
						return err
					}

					site = &sites[selected]
				}
			default:
				site, err = cfg.FindSiteByHostName(siteArg)
				if err != nil {
					return err
				}
			}

			output.Info("Setting permissions for", site.Hostname)

			// find the docker executable
			cli, err := exec.LookPath("docker")
			if err != nil {
				return err
			}

			for _, f := range files {
				stat, _ := docker.ContainerStatPath(ctx, site.Hostname, f.path)

				if stat.Name != "" {
					output.Pending("modifying", f.path)

					containerUser := "www-data"
					if runtime.GOOS == "linux" {
						user, err := user.Current()
						if err != nil {
							return err
						}
						containerUser = fmt.Sprintf("%s:%s", user.Uid, user.Gid)
					}

					cmds := []string{"exec", "-it", "--user", containerUser, site.Hostname}

					if f.directory {
						cmds = append(cmds, "chmod", "-R", "777", f.path)
					} else {
						cmds = append(cmds, "chmod", "777", f.path)
					}

					// create the command
					c := exec.Command(cli, cmds...)

					c.Stdin = cmd.InOrStdin()
					c.Stderr = cmd.ErrOrStderr()
					c.Stdout = cmd.OutOrStdout()

					if err := c.Run(); err != nil {
						fmt.Println("error setting permissions on", f.path)
					}

					output.Done()
				}
			}

			return nil
		},
	}

	return cmd
}
