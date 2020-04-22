package helpers

import (
	"errors"
	"os"
	"path/filepath"
)

// FindWebRoot takes a directory and will search for the "webroot" automatically.
// if it cannot find a know webroot, the func will return an error. This is
// used when determining a sites complete path to the webroot for nginx.
func FindWebRoot(path string) (string, error) {
	var w string
	if err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		if !info.IsDir() {
			return nil
		}

		switch info.Name() {
		case "web":
			w = info.Name()
		case "public":
			w = info.Name()
		case "public_html":
			w = info.Name()
		case "www":
			w = info.Name()
		}

		return nil
	}); err != nil {
		return "", err
	}

	if w == "" {
		return "", errors.New("unable to locate the webroot for " + path)
	}

	return w, nil
}
