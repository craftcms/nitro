package create

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command/create/internal/urlgen"
	"github.com/craftcms/nitro/pkg/downloader"
	"github.com/craftcms/nitro/pkg/envedit"
	"github.com/craftcms/nitro/pkg/pathexists"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
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

				dir = filepath.Join(args[1])
			default:
				// only the directory was provided, download craft to that directory
				u, err := urlgen.Generate("")
				if err != nil {
					return err
				}

				download = u

				dir = filepath.Join(args[0])
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
			_, err := prompt.CreateSite(home, dir, output)
			if err != nil {
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
					key := uuid.New()

					// update the env
					update, err := envedit.Edit(envFilePath, map[string]string{
						"SECURITY_KEY": key.String(),
						"DB_SERVER":    dbhost,
						"DB_DATABASE":  dbname,
						"DB_PORT":      port,
						"DB_DRIVER":    driver,
						"DB_USER":      "nitro",
						"DB_PASSWORD":  "nitro",
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

	return cmd
}
