package add

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/envedit"
	"github.com/craftcms/nitro/pkg/pathexists"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # add the current directory as a site
  nitro add

  # add a directory as the site
  nitro add my-project`

// NewCommand returns the command to add a site to the nitro config.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Short:   "Adds a site.",
		Example: exampleText,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.VerifyInit(cmd, args, home, output)
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return prompt.RunApply(cmd, args, false, output)
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
				// check if the path is using the ~
				if strings.HasPrefix(args[0], "~") {
					dir = strings.Replace(args[0], "~", home, 1)
				} else {
					dir = filepath.Join(args[0])
				}

				// make sure the directory exists
				if !pathexists.IsDirectory(dir) {
					return fmt.Errorf("unable to find the directory: %s", dir)
				}
			default:
				dir = filepath.Clean(wd)
			}

			output.Info("Adding site…")

			if _, err := prompt.CreateSite(home, dir, output); err != nil {
				return err
			}

			exampleEnv := filepath.Join(dir, ".env.example")
			envFilePath := filepath.Join(dir, ".env")

			// check if the directory has a .env-example
			if pathexists.IsFile(exampleEnv) && !pathexists.IsFile(envFilePath) {
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
			database, dbhost, dbname, port, driver, err := prompt.CreateDatabase(cmd, docker, output)
			if err != nil {
				return err
			}

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

					// truncate the file
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

			output.Info("New site added! 🎉")

			return nil
		},
	}

	return cmd
}
