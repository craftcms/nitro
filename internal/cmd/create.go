package cmd

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pixelandtonic/prompt"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/helpers"
)

var createcommand = &cobra.Command{
	Use:   "create",
	Short: "Create a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = flagMachineName
		p := prompt.NewPrompt()

		dir, err := p.Ask("What is the name of the project?", &prompt.InputOptions{
			AppendQuestionMark: false,
		})
		if err != nil {
			return err
		}

		if helpers.DirExists(dir) {
			return errors.New(fmt.Sprintf("the directory %q already exists", dir))
		}

		// download craft
		url := "https://craftcms.com/latest-v3.zip"
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Write the body to file
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// unzip the folder
		zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
		if err != nil {
			return err
		}

		if err := os.Mkdir(dir, 0744); err != nil {
			return err
		}

		for _, f := range zr.File {
			p := filepath.Join(dir, f.Name)

			// if this is a directory, lets make it
			if f.FileInfo().IsDir() {
				if err := os.MkdirAll(p, os.ModePerm); err != nil {
					return err
				}
			}

			// if its a file, make it
			file, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			rdr, err := f.Open()
			if err != nil {
				return err
			}

			if _, err := io.Copy(file, rdr); err != nil {
				return err
			}
			if err := file.Close(); err != nil {
				return err
			}

			if err := rdr.Close(); err != nil {
				return err
			}
		}

		return nil
	},
}
