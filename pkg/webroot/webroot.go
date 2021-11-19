package webroot

import (
	"fmt"
	"io/ioutil"
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
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	// check the directories
	for _, info := range dir {
		// only check directories
		if !info.IsDir() {
			continue
		}

		// is this a known webroot?
		if info.Name() == "web" || info.Name() == "public" || info.Name() == "public_html" || info.Name() == "html" {
			return info.Name(), nil
		}
	}

	// always return web if we can't find anything here
	return "web", nil
}
