package create

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command/create/internal/urlgen"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # create a new default craft project (similar to "composer create-project craftcms/craft my-project")
  nitro create my-project

  # bring your own git repo
  nitro create https://github.com/craftcms/demo my-project

  # you can also provide shorthand urls for github
  nitro create craftcms/demo my-project`

var download = "https://github.com/craftcms/craft/archive/HEAD.zip"

// New returns the create command to automate the process of setting up a new Craft project.
// It also allows you to pass an option argument that is a URL to your own github repo.
func New(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create project",
		Example: exampleText,
		Args:    cobra.MinimumNArgs(1),
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
				// only the directory was provided, download craft to that repo
				u, err := urlgen.Generate("")
				if err != nil {
					return err
				}

				download = u
				dir = cleanDirectory(args[0])
			}

			output.Pending("setting up project")

			// download the zip
			resp, err := http.Get(download.String())
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("\nUnable to download %s.\nStatus: %d", download.String(), resp.StatusCode)
			}

			// create a temp file
			file, err := ioutil.TempFile(os.TempDir(), "nitro-create-download-")
			if err != nil {
				return err
			}
			defer file.Close()
			defer os.Remove(file.Name())

			// copy the download into the new file
			_, err = io.Copy(file, resp.Body)
			if err != nil {
				return fmt.Errorf("unable to copy the file, %w", err)
			}

			// extract the zip
			r, err := zip.OpenReader(file.Name())
			if err != nil {
				return err
			}
			defer r.Close()

			// TODO(jasonmccallister) ask for the version of PHP

			for _, f := range r.File {
				// github archives has a nested folder, so we need to trim the first directory
				p := strings.Split(f.Name, string(os.PathSeparator))
				fpath := filepath.Join(dir, strings.Join(p[1:], string(os.PathSeparator)))

				// if !strings.HasPrefix(fpath, filepath.Clean(dir)+string(os.PathSeparator)) {
				// 	return fmt.Errorf("%s: illegal file path", fpath)
				// }

				if f.FileInfo().IsDir() {
					os.MkdirAll(fpath, os.ModePerm)

					continue
				}

				if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
					return err
				}

				out, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
				if err != nil {
					return err
				}

				rc, err := f.Open()
				if err != nil {
					return err
				}

				_, err = io.Copy(out, rc)

				// Close the file without defer to close before next iteration of loop
				out.Close()
				rc.Close()
			}

			output.Done()

			output.Info("project created 🤓")

			// TODO(jasonmccallister) prompt the user for the version of php, webroot, hostname
			// prompt for the php version
			// versions := phpversions.Versions
			// selected, err := output.Select(cmd.InOrStdin(), "Choose a PHP version: ", versions)
			// if err != nil {
			// 	return err
			// }

			// run the composer install command
			for _, c := range cmd.Parent().Commands() {
				if c.Use == "composer" {
					// run composer install using the new directory
					// we pass the command itself instead of the parent
					// command
					if err := c.RunE(c, []string{dir, "--version=" + cmd.Flag("composer-version").Value.String()}); err != nil {
						return err
					}
				}
			}

			// TODO(jasonmccallister) edit the .env
			// TODO(jasonmccallister) ask if we should run apply now

			return nil
		},
	}

	// TODO(jasonmccallister) add flags for the composer and node versions
	cmd.Flags().String("composer-version", "2", "version of composer to use")
	cmd.Flags().String("node-version", "14", "version of node to use")

	return cmd
}

func cleanDirectory(s string) string {
	return strings.TrimSpace(strings.Replace(s, " ", "-", -1))
}
