package permissions

import (
	"fmt"
	"os/exec"
	"os/user"
	"runtime"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/terminal"
)

type file struct {
	path      string
	directory bool
}

const exampleText = `  # fix filesystem permissions for craft
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

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "permissions",
		Short:   "Fix Craft permissions",
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

			// find the docker executable
			cli, err := exec.LookPath("docker")
			if err != nil {
				return err
			}

			for _, f := range files {
				stat, _ := docker.ContainerStatPath(ctx, "tutorial.nitro", f.path)

				if stat.Name != "" {
					fmt.Println("setting permissions on", f.path)

					containerUser := "www-data"
					if runtime.GOOS == "linux" {
						user, err := user.Current()
						if err != nil {
							return err
						}
						containerUser = fmt.Sprintf("%s:%s", user.Uid, user.Gid)
					}

					cmds := []string{"exec", "-it", "--user", containerUser, "tutorial.nitro"}

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
				}
			}

			return nil
		},
	}

	return cmd
}
