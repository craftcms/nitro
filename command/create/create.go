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

	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var (
	// ErrExample is used when we want to share an error
	ErrExample = fmt.Errorf("some example error")
)

const exampleText = `  # create command
  # create a new default craft project
  nitro create

  # bring your own git repo
  nitro create https://github.com/craftcms/demo`

var download = "https://github.com/craftcms/craft/archive/HEAD.zip"

// New returns the create command to automate the process of setting up a new Craft project.
// It also allows you to pass an option argument that is a URL to your own github repo.
func New(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create project",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// get the url from args or the default
			var u url.URL
			if len(args) > 0 {
				parsed, err := url.Parse(args[0])
				if err != nil {
					return fmt.Errorf("")
				}
				u = *parsed
			} else {
				parsed, err := url.Parse(download)
				if err != nil {
					return fmt.Errorf("")
				}
				u = *parsed
			}

			// dir := "docker/"

			// https://github.com/craftcms/craft/archive/HEAD.zip
			output.Pending("setting up project")

			// download the zip
			resp, err := http.Get(u.String())
			if err != nil {
				return err
			}
			defer resp.Body.Close()

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

			// TODO(jasonmccallister) make this dynamic
			dir := "docker"
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

			output.Info("new project created 🤓")

			// TODO(jasonmccallister) run the composer install command
			// composerCommand := composer.New(docker, output)
			// composerCommand.Flags().Set("version", "1")
			// composerCommand.RunE(, []string{dir})

			// TODO(jasonmccallister) edit the .env
			// TODO(jasonmccallister) ask if we should run apply now

			return nil
		},
	}

	return cmd
}
