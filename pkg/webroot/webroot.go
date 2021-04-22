package webroot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ErrNotFound is returned when unable to find a web root for a specified path
	ErrNotFound = fmt.Errorf("unable to locate a web root")
)

// Find takes a path and will check for the web root of the
// project. Find will look for web, public, and public_html
// directories and return when the first directory match.
// If it cannot find the web root it will return an error.
func Find(path string) (string, error) {
	var root string
	if err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		// safety check
		if info == nil {
			return nil
		}

		// ignore files
		if !info.IsDir() {
			return nil
		}

		// check if the directory is relative to the vendor dir
		dir, e := filepath.Rel(filepath.Join(path), filepath.Join(p))
		if e != nil {
			return nil
		}

		// if the dir is in the vendor dir, we want to ignore it
		if strings.Contains(dir, fmt.Sprintf("vendor%c", os.PathSeparator)) || strings.Contains(dir, fmt.Sprintf("node_modules%c", os.PathSeparator)) {
			return nil
		}

		// if the directory is considered a web root
		if info.Name() == "web" || info.Name() == "public" || info.Name() == "public_html" || info.Name() == "html" {
			root = info.Name()
		}

		// if its not set, keep trying
		if root != "" {
			return nil
		}

		return nil
	}); err != nil {
		return "", err
	}

	// if we found the root, return it
	if root != "" {
		return root, nil
	}

	return "", ErrNotFound
}
