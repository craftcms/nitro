package downloader

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Getter is an interface for getting the contents of a url
// and a directory to unzip the contents into a directory.
type Getter interface {
	Get(url, dir string) error
}

// Downloader wraps the HTTP client to make get requests to a url
// and unzip to a specific directory.
type Downloader struct {
	client *http.Client
}

// NewDownloader creates a downloader with a default HTTP client
// which is used to download files from the net.
func NewDownloader() *Downloader {
	return &Downloader{
		client: http.DefaultClient,
	}
}

// Get takes a url and a directory where the contents should be
// unzipped into.
func (d *Downloader) Get(url, dir string) error {
	// download the zip
	resp, err := d.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// check the response code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to download %s. status: %d", url, resp.StatusCode)
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

	if err := unzip(file, dir); err != nil {
		return err
	}

	return nil
}

func unzip(file *os.File, dir string) error {
	// extract the zip
	r, err := zip.OpenReader(file.Name())
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// github archives has a nested folder, so we need to trim the first directory
		p := strings.Split(f.Name, string(os.PathSeparator))
		fpath := filepath.Join(dir, strings.Join(p[1:], string(os.PathSeparator)))

		// if !strings.HasPrefix(fpath, filepath.Clean(dir)+string(os.PathSeparator)) {
		// 	return fmt.Errorf("%s: illegal file path", fpath)
		// }

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}

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
		if err != nil {
			return err
		}

		// Close the file without defer to close before next iteration of loop
		out.Close()
		rc.Close()
	}

	return nil
}
