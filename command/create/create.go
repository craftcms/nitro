package create

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command/create/internal/urlgen"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/downloader"
	"github.com/craftcms/nitro/pkg/envedit"
	"github.com/craftcms/nitro/pkg/pathexists"
	"github.com/craftcms/nitro/pkg/phpversions"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/craftcms/nitro/pkg/validate"
	"github.com/craftcms/nitro/pkg/webroot"
)

const exampleText = `  # create a new default craft project (similar to "composer create-project craftcms/craft my-project")
  nitro create my-project

  # bring your own git repo
  nitro create https://github.com/craftcms/demo my-project

  # you can also provide shorthand urls for github
  nitro create craftcms/demo my-project`

// NewCommand returns the create command to automate the process of setting up a new Craft project.
// It also allows you to pass an option argument that is a URL to your own github repo.
func NewCommand(home string, docker client.CommonAPIClient, getter downloader.Getter, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create project",
		Example: exampleText,
		Args:    cobra.MinimumNArgs(1),
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
			// get the url from args or the default
			var download *url.URL
			var dir string

			switch len(args) {
			case 2:
				// the directory and url are specified
				u, err := urlgen.Generate(args[0])
				if err != nil {
					return err
				}

				download = u
				dir = cleanDirectory(args[1])
			default:
				// only the directory was provided, download craft to that directory
				u, err := urlgen.Generate("")
				if err != nil {
					return err
				}

				download = u
				dir = cleanDirectory(args[0])
			}

			// check if the directory already exists
			if exists := pathexists.IsDirectory(dir); exists {
				return fmt.Errorf("directory %q already exists", dir)
			}

			output.Info("Downloading", download.String(), "...")

			output.Pending("setting up project")

			// download the file
			if err := getter.Get(download.String(), dir); err != nil {
				return err
			}

			output.Done()

			output.Info("New site downloaded 🤓")

			// --- done with download

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

			// walk the user through the site
			if err := promptSiteAdd(home, dir, output); err != nil {
				return err
			}

			//  prompt for a new database
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

			// run the composer install command
			for _, c := range cmd.Parent().Commands() {
				if c.Use == "composer" {
					// ask the user if we should install composer dependencies
					installComposerDeps, err := output.Confirm("Should we install composer dependencies?", false, "")
					if err != nil {
						return err
					}

					if installComposerDeps {
						// change into the projects new directory for the composer install
						if err := os.Chdir(filepath.Join(dir)); err != nil {
							break
						}

						// run composer install using the new directory
						// we pass the command itself instead of the parent
						// command
						if err := c.RunE(c, []string{"install", "--ignore-platform-reqs"}); err != nil {
							output.Info(err.Error())
							break
						}
					}
				}
			}

			return nil
		},
	}

	// TODO(jasonmccallister) add flags for the composer and node versions
	cmd.Flags().String("composer-version", "2", "version of composer to use")
	cmd.Flags().String("node-version", "14", "version of node to use")

	return cmd
}

func cleanDirectory(s string) string {
	return filepath.Join(s)
}

func promptSiteAdd(home, dir string, output terminal.Outputer) error {
	// create a new site
	site := config.Site{}

	// get the hostname from the directory
	// p := filepath.Join(dir)
	sp := strings.Split(filepath.Join(dir), string(os.PathSeparator))
	site.Hostname = sp[len(sp)-1]

	// append the test domain if there are no periods
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
	siteAbsPath, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	site.Path = strings.Replace(siteAbsPath, home, "~", 1)

	output.Success("adding site", site.Path)

	// get the web directory
	found, err := webroot.Find(dir)
	if err != nil {
		return err
	}

	// if the root is still empty, we fall back to the default
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
	selected, err := output.Select(os.Stdin, "Choose a PHP version: ", versions)
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

	return nil
}
