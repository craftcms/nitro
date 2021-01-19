package add

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/phpversions"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/craftcms/nitro/pkg/validate"
	"github.com/craftcms/nitro/pkg/webroot"
)

// TODO - prompt user for the database engine and new database
// TODO - edit the env file for the user

const exampleText = `  # add the current project as a site
  nitro add

  # add a directory as the site
  nitro add my-project`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Short:   "Add a site",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			output.Info("Adding site…")

			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// get working directory or provided arg
			var dir string
			switch len(args) {
			case 1:
				dir = filepath.Join(wd, args[0])
			default:
				dir = filepath.Clean(wd)
			}

			// create a new site
			site := config.Site{}

			// get the hostname from the directory
			sp := strings.Split(dir, string(os.PathSeparator))
			site.Hostname = sp[len(sp)-1]

			// append the default TLD if there are no periods in the path name
			if !strings.Contains(site.Hostname, ".") {
				// set the default tld
				tld := "nitro"
				if os.Getenv("NITRO_DEFAULT_TLD") != "" {
					tld = os.Getenv("NITRO_DEFAULT_TLD")
				}

				site.Hostname = fmt.Sprintf("%s.%s", site.Hostname, tld)
			}

			// prompt for the hostname
			hostname, err := output.Ask("Enter the hostname", site.Hostname, ":", &validate.HostnameValidator{})
			if err != nil {
				return err
			}

			// set the input as the hostname
			site.Hostname = hostname

			output.Success("setting hostname to", site.Hostname)

			// set the sites directory but make the path relative
			site.Path = strings.Replace(dir, home, "~", 1)

			output.Success("adding site", site.Path)

			// get the web directory
			found, err := webroot.Find(dir)
			if err != nil {
				return err
			}

			if found == "" {
				found = "web"
			}

			// set the webroot
			site.Webroot = found

			// prompt for the webroot
			root, err := output.Ask("Enter the webroot for the site", site.Webroot, ":", nil)
			if err != nil {
				return err
			}

			site.Webroot = root

			output.Success("using webroot", site.Webroot)

			// prompt for the php version
			versions := phpversions.Versions
			selected, err := output.Select(cmd.InOrStdin(), "Choose a PHP version: ", versions)
			if err != nil {
				return err
			}

			// set the version of php
			site.Version = versions[selected]

			output.Success("setting PHP version", site.Version)

			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// add the site to the config
			if err := cfg.AddSite(site); err != nil {
				return err
			}

			output.Pending("saving file")

			// save the config file
			if err := cfg.Save(); err != nil {
				output.Warning()

				return err
			}

			output.Done()

			output.Info("Site added 🌍")

			// ask if the apply command should run
			fmt.Print("Apply changes now [Y/n]? ")

			s := bufio.NewScanner(os.Stdin)
			s.Split(bufio.ScanLines)

			confirm := true
			for s.Scan() {
				txt := s.Text()

				txt = strings.TrimSpace(txt)

				if txt == "" {
					confirm = true
					break
				}

				for _, answer := range []string{"n", "N", "no", "No", "NO"} {
					if txt == answer {
						confirm = false
						break
					}
				}
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

	return cmd
}
