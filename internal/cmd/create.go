package cmd

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pixelandtonic/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/internal/config"
	"github.com/craftcms/nitro/internal/helpers"
	"github.com/craftcms/nitro/internal/validate"
)

const craftDownloadURL = "https://craftcms.com/latest-v3.zip"

var createcommand = &cobra.Command{
	Use:   "create",
	Short: "Create a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = flagMachineName
		p := prompt.NewPrompt()
		// load the config
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		// prompt the user for input
		dir, err := p.Ask("What is the name of the project?", &prompt.InputOptions{
			AppendQuestionMark: false,
		})
		if err != nil {
			return err
		}

		// clean up the directory name
		dir = strings.TrimSpace(strings.Replace(dir, " ", "-", -1))

		// check if the directory exists
		if helpers.DirExists(dir) {
			return errors.New(fmt.Sprintf("Directory %q already exists", dir))
		}

		// prompt the user for the sites hostname, default to the directory name
		hostname, err := p.Ask("Enter the hostname", &prompt.InputOptions{
			Default:   dir,
			Validator: validate.Hostname,
		})
		if err != nil {
			return err
		}

		// download craft
		file, err := download()
		if err != nil {
			return err
		}

		projectPath, err := filepath.Abs(dir)
		if err != nil {
			return err
		}

		// unzip the files
		if err := unzip(file, projectPath); err != nil {
			return err
		}

		// apply flags to add command
		args = append(args, dir)
		add := addCommand
		if err := add.Flag("hostname").Value.Set(hostname); err != nil {
			return err
		}
		if err := add.Flag("webroot").Value.Set("web"); err != nil {
			return err
		}

		return add.RunE(cmd, args)
	},
}

func download() (*os.File, error) {
	r, err := http.Get(craftDownloadURL)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	// create the temp file
	f, err := ioutil.TempFile(os.TempDir(), "nitro-craft-cms-download")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fmt.Println("Downloading Craft CMS...")

	// copy the downloaded contents into the temp file
	_, err = io.Copy(f, r.Body)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func unzip(source *os.File, path string) error {
	r, err := zip.OpenReader(source.Name())
	if err != nil {
		return err
	}
	defer r.Close()

	// move files into place
	for _, f := range r.File {
		path := filepath.Join(path, f.Name)

		// zipslip
		//if !strings.HasPrefix(path, filepath.Clean(path)+string(os.PathSeparator)) {
		//	log.Printf("%s: illegal file path", path)
		//	continue
		//}

		// create if a directory
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		// create the files
		if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}

		out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		readerCloser, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(out, readerCloser)

		// cleanup
		out.Close()
		readerCloser.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
