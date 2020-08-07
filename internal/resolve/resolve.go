package resolve

import (
	"os"
	"path/filepath"
	"strings"
)

// AbsPath takes a path and the users home directory
// and will resolve the path with filepath.Abs
// this takes into account the ~ home dir references
func AbsPath(path, home string) (string, error) {
	// if this is a relative path, replace with the users home dir
	if strings.HasPrefix(path, "~/") {
		path = strings.Replace(path, "~", home, 1)
	}

	// get the absolute path
	p, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// check if the file exists
	if _, err := os.Stat(p); err != nil {
		return "", err
	}

	return p, nil
}
