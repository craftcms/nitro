package helpers

import (
	"errors"
	"os"
	"path/filepath"
)

func FindWebRoot(path string) (string, error) {
	var webroot string
	if err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		if !info.IsDir() {
			return nil
		}

		switch info.Name() {
		case "web":
			webroot = info.Name()
		case "public":
			webroot = info.Name()
		case "public_html":
			webroot = info.Name()
		case "www":
			webroot = info.Name()
		}

		return nil
	}); err != nil {
		return "", err
	}

	if webroot == "" {
		return "", errors.New("unable to locate the webroot")
	}

	return webroot, nil
}
