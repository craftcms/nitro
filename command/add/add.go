package add

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var (
	// ErrExample is used when we want to share an error
	ErrExample = fmt.Errorf("some example error")
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
			env := cmd.Flag("environment").Value.String()
			output.Info("Adding site...")

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

			// append the test domain if there are no periods
			if strings.Contains(site.Hostname, ".") == false {
				// set the default tld
				tld := "nitro"
				if os.Getenv("NITRO_DEFAULT_TLD") != "" {
					tld = os.Getenv("NITRO_DEFAULT_TLD")
				}

				site.Hostname = fmt.Sprintf("%s.%s", site.Hostname, tld)
			}

			// prompt for the hostname
			fmt.Print(fmt.Sprintf("Enter the hostname [%s]: ", site.Hostname))
			for {
				rdr := bufio.NewReader(os.Stdin)
				char, _ := rdr.ReadString('\n')

				// remove the carriage return
				char = strings.TrimRight(char, "\n")

				// does it have spaces?
				if strings.ContainsAny(char, " ") {
					fmt.Println("Please enter a hostname without spaces 🙄...")
					fmt.Print(fmt.Sprintf("Enter the hostname [%s]: ", site.Hostname))

					continue
				}

				// if its empty, we are setting the default
				if char == "" {
					break
				}

				// set the input as the hostname
				site.Hostname = char
			}

			output.Success("setting hostname to", site.Hostname)

			// set the sites directory but make the path relative
			site.Path = strings.Replace(dir, home, "~", 1)

			output.Success("adding site", site.Path)

			// get the web directory
			var root string
			if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				// don't go into subdirectories and ignore files
				if path != dir || info.IsDir() == false {
					return nil
				}

				// if the directory is considered a web root
				if info.Name() == "web" || info.Name() == "public" || info.Name() == "public_html" {
					root = info.Name()
				}

				// if its not set, keep trying
				if root != "" {
					return nil
				}

				return nil
			}); err != nil {
				return err
			}

			// if the root is still empty, we fall back to the default
			if root == "" {
				root = "web"
			}

			// set the webroot
			site.Dir = root

			// prompt for the webroot
			fmt.Print(fmt.Sprintf("Enter the webroot for the site [%s]: ", site.Dir))
			for {
				rdr := bufio.NewReader(os.Stdin)
				input, _ := rdr.ReadString('\n')

				// remove the carriage return
				input = strings.TrimRight(input, "\n")

				// does it have spaces?
				if strings.ContainsAny(input, " ") {
					fmt.Println("Please enter a webroot without spaces 🙄...")
					fmt.Print(fmt.Sprintf("Enter the webroot for the site [%s]: ", site.Dir))

					continue
				}

				// if its empty, we are setting the default
				if input == "" {
					break
				}

				// set the input as the hostname
				site.Dir = input
				break
			}

			output.Success("using webroot", site.Dir)

			// prompt for the php version
			versions := []string{"7.4", "7.3", "7.2", "7.1"}
			selected, err := output.Select(cmd.InOrStdin(), "Choose a PHP version: ", versions)
			if err != nil {
				return err
			}

			// set the version of php
			site.Version = versions[selected]

			output.Success("setting PHP version", site.Version)

			// load the config
			cfg, err := config.Load(home, env)
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

			output.Info(fmt.Sprintf("Site added to %s 🌍", env))

			// ask if the apply command should run
			var response string
			fmt.Print("Apply changes now [Y/n]? ")
			if _, err := fmt.Scanln(&response); err != nil {
				return fmt.Errorf("unable to provide a prompt, %w", err)
			}

			// get the response
			resp := strings.TrimSpace(response)
			var confirm bool
			for _, answer := range []string{"y", "Y", "yes", "Yes", "YES"} {
				if resp == answer {
					confirm = true
				}
			}

			// we are skipping the apply step
			if confirm == false {
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
