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

const exampleText = `  # add the current project as a site
  nitro add

  # add a directory as the site
  nitro add my-project`

// New is used for scaffolding new commands
func New(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Short:   "Add a new site",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			output.Info("Adding site...")

			// get the environment
			site := config.Site{}

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

			// get the hostname from the directory
			sp := strings.Split(dir, string(os.PathSeparator))
			site.Hostname = sp[len(sp)-1]

			// append the test domain if there are no periods
			if strings.Contains(site.Hostname, ".") == false {
				// set the default tld
				tld := "test"
				if os.Getenv("NITRO_DEFAULT_TLD") != "" {
					tld = os.Getenv("NITRO_DEFAULT_TLD")
				}

				site.Hostname = fmt.Sprintf("%s.%s", site.Hostname, tld)
			}

			// prompt for the hostname
			fmt.Print(fmt.Sprintf("Enter the hostname [%s]: ", site.Hostname))
			w := true
			for w {
				rdr := bufio.NewReader(os.Stdin)
				char, _ := rdr.ReadString('\n')

				// remove the carriage return
				char = strings.TrimRight(char, "\n")

				// does it have spaces?
				if strings.ContainsAny(char, " ") {
					w = true
					fmt.Println("Please enter a hostname without spaces 🙄...")
					fmt.Print(fmt.Sprintf("Enter the hostname [%s]: ", site.Hostname))
				} else {
					site.Hostname = char
					w = false
				}
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

					return nil
				}

				return nil
			}); err != nil {
				return err
			}

			if root == "" {
				root = "web"
			}

			// set the webroot
			site.Dir = root

			output.Success("using webroot", site.Dir)

			// prompt for the php version
			versions := []string{"7.4", "7.3", "7.2", "7.1"}
			selected, err := output.Select(cmd.InOrStdin(), "Choose a PHP version: ", versions)
			if err != nil {
				return err
			}

			// set the version of php
			site.PHP = versions[selected]

			output.Success("setting PHP version", site.PHP)

			fmt.Println(site)

			// verify the site does not exist

			// add the site to the config
			return nil
		},
	}

	return cmd
}
