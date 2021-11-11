package add

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/paths"
	"github.com/craftcms/nitro/pkg/validate"
	"github.com/craftcms/nitro/pkg/webroot"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/envedit"
	"github.com/craftcms/nitro/pkg/pathexists"
	"github.com/craftcms/nitro/pkg/prompt"
	"github.com/craftcms/nitro/pkg/terminal"
)

var flagExcludeDependencies bool

// NewCommand returns the command to add an app to the nitro config.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Adds an app.",
		Example: `  # add the current working directory as an app
  nitro add

  # add the specified directory as the app directory
  nitro add my-project`,
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
				// check if the arg is using a relative path ~
				if strings.HasPrefix(args[0], "~") {
					dir, err = paths.Clean(home, args[0])
					if err != nil {
						return fmt.Errorf("unable to locate the directory, err: %s", err)
					}
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

			output.Info("Adding app…")

			app := config.App{}

			// are we excluding any dependencies?
			if flagExcludeDependencies {
				app.Excludes = []string{"node_modules", "vendor"}
			}

			// generate a hostname for the app using the directory
			sp := strings.Split(filepath.Join(dir), string(os.PathSeparator))
			app.Hostname = sp[len(sp)-1]
			// append the nitro domain if there are no periods in the hostname
			if !strings.Contains(app.Hostname, ".") {
				// set the default tld
				tld := "nitro"
				if os.Getenv("NITRO_DEFAULT_TLD") != "" {
					tld = os.Getenv("NITRO_DEFAULT_TLD")
				}

				app.Hostname = fmt.Sprintf("%s.%s", app.Hostname, tld)
			}

			// prompt the user to validate the hostname
			hostname, err := output.Ask("Enter the hostname", app.Hostname, ":", &validate.HostnameValidator{})
			if err != nil {
				return err
			}
			output.Success("setting the app hostname to", hostname)

			// set the apps path and replace the full path with a relative using ~
			abs, err := filepath.Abs(dir)
			if err != nil {
				return fmt.Errorf("unable to find the absolute path to the app, err: %s", err)
			}
			app.Path = strings.Replace(abs, home, "~", 1)
			output.Success("adding app", app.Path)

			// find the apps webroot from the directory
			found, err := webroot.Find(dir)
			if err != nil {
				return fmt.Errorf("unable to find the webroot for the app, err: %s", err)
			}

			// prompt for the web root
			root, err := output.Ask("Enter the web root for the app", found, ":", nil)
			if err != nil {
				return err
			}
			app.Webroot = root
			output.Success("using web root", app.Webroot)

			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// add the app
			if err := cfg.AddApp(app); err != nil {
				return err
			}

			// save the config
			if err := cfg.Save(); err != nil {
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

			// always set default environment variables
			envVars := map[string]string{
				"DB_USER":     "nitro",
				"DB_PASSWORD": "nitro",
			}

			// if the user selected a database, add that information
			if database {
				envVars["DB_SERVER"] = dbhost
				envVars["DB_PORT"] = port
				envVars["DB_DATABASE"] = dbname
				envVars["DB_DRIVER"] = driver
			}

			// if they wanted a new database edit the env
			if pathexists.IsFile(envFilePath) {
				// ask the user if we should update the .env?
				updateEnv, err := output.Confirm("Should we update the env file?", false, "")
				if err != nil {
					return err
				}

				// the user wants to update the env file
				if updateEnv {
					// check if the security key is already set
					if !envedit.EnvExists(envFilePath, "SECURITY_KEY") {
						envVars["SECURITY_KEY"] = uuid.New().String()
					}

					// update the env
					update, err := envedit.Edit(envFilePath, envVars)
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

			output.Info("New app", app.Hostname, "added! 🎉")

			return nil
		},
	}

	cmd.Flags().BoolVar(&flagExcludeDependencies, "exclude-deps", false, "Ignore node_modules and vendor directories when creating the app container.")

	return cmd
}
