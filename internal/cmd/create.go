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

	"github.com/craftcms/nitro/internal/helpers"
)

const craftDownloadURL = "https://craftcms.com/latest-v3.zip"

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

		// clean up the directory name
		dir = strings.TrimSpace(strings.Replace(dir, " ", "-", -1))

		if helpers.DirExists(dir) {
			return errors.New(fmt.Sprintf("the directory %q already exists", dir))
		}

		projectPath, err := filepath.Abs(dir)
		if err != nil {
			return err
		}

		// download craft
		resp, err := http.Get(craftDownloadURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		tmpDir := os.TempDir()

		// create the temp file
		tempFile, err := ioutil.TempFile(tmpDir, "nitro")
		if err != nil {
			return err
		}
		defer tempFile.Close()

		fmt.Println("Downloading Craft CMS...")

		// copy the downloaded contents into the temp file
		_, err = io.Copy(tempFile, resp.Body)
		if err != nil {
			return err
		}

		// unzip the files
		if err := unzip(tempFile, projectPath); err != nil {
			return err
		}

		return nil
	},
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
		if !strings.HasPrefix(path, filepath.Clean(path)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", path)
		}

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
