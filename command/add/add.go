package add

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/envedit"
	"github.com/craftcms/nitro/pkg/pathexists"
	"github.com/craftcms/nitro/pkg/phpversions"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/craftcms/nitro/pkg/validate"
	"github.com/craftcms/nitro/pkg/webroot"
)

const exampleText = `  # add the current project as a site
  nitro add

  # add a directory as the site
  nitro add my-project`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Short:   "Add a site",
		Example: exampleText,
		PostRunE: func(cmd *cobra.Command, args []string) error {
			// ask if the apply command should run
			apply, err := output.Confirm("Apply changes now", true, "?")
			if err != nil {
				return err
			}

			// if apply is false return nil
			if !apply {
				return nil
			}

			// run the apply command
			for _, c := range cmd.Parent().Commands() {
				// set the apply command
				if c.Use == "apply" {
					return c.RunE(c, args)
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
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
				// make sure the directory exists
				if exists := pathexists.IsDirectory(dir); !exists {
					return fmt.Errorf("unable to find the directory: %s", dir)
				}
			default:
				dir = filepath.Clean(wd)
			}

			output.Info("Adding site…")

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

			exampleEnv := filepath.Join(dir, ".env.example")

			// check if the directory has a .env-example
			if pathexists.IsFile(exampleEnv) {
				// open the example
				example, err := os.Open(exampleEnv)
				if err != nil {
					output.Info("unable to open the file", exampleEnv)
				}
				defer example.Close()

				// create the env file
				env, err := os.Create(filepath.Join(dir, ".env"))
				if err != nil {
					output.Info("unable to create the file", filepath.Join(dir, ".env"))
				}
				defer env.Close()

				if _, err := io.Copy(env, example); err != nil {
					output.Info("unable to copy the example env")
				}
			}

			// prompt for a database
			database, dbhost, dbname, port, driver, err := prompt.CreateDatabase(cmd.Context(), docker, output)
			if err != nil {
				return err
			}

			envFilePath := filepath.Join(dir, ".env")

			// if the wanted a new database edit the env
			if database && pathexists.IsFile(envFilePath) {
				// ask the user if we should update the .env?
				updateEnv, err := output.Confirm("Should we update the env file?", false, "")
				if err != nil {
					return err
				}

				if updateEnv {
					// update the env
					update, err := envedit.Edit(envFilePath, map[string]string{
						"DB_SERVER":   dbhost,
						"DB_DATABASE": dbname,
						"DB_PORT":     port,
						"DB_DRIVER":   driver,
						"DB_USER":     "nitro",
						"DB_PASSWORD": "nitro",
					})
					if err != nil {
						output.Info("unable to edit the env")
					}

					// open the file
					f, err := os.OpenFile(envFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
					if err != nil {
						return err
					}
					defer f.Close()

					// trunacte the file
					if err := f.Truncate(0); err != nil {
						return err
					}

					// write the new contents
					if _, err := f.Write([]byte(update)); err != nil {
						return err
					}

					output.Info(".env updated!")
				}
			}

			return nil
		},
	}

	return cmd
}
