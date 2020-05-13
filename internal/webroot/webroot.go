package webroot

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Find takes a directory and will search for the "webroot" automatically.
// if it cannot find a know webroot, the func will return an error. This is
// used when determining a sites complete path to the webroot for nginx.
func Find(path string) (string, error) {
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

// Matches takes the found webroot from a nitro machine
// and compares the config webroot and returns
// a boolean if they match.
func Matches(output, webroot string) (bool, string) {
	if len(output) == 0 {
		return false, ""
	}

	sp := strings.Split(strings.TrimSpace(output), " ")

	// remove the trailing ;
	sp[1] = strings.TrimRight(sp[1], ";")

	if sp[1] == webroot {
		return true, sp[1]
	}

	return false, sp[1]
}
