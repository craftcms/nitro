package paths

import (
	"path/filepath"
	"strings"
)

// Clean takes a users home directory and a path and returns the complete
// path to the provided path.
func Clean(home, path string) (string, error) {
	p := path
	if strings.Contains(p, "~") {
		p = strings.Replace(p, "~", home, -1)
	}

	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}

	return filepath.Clean(abs), nil
}

// MakeRelative takes a users home directory and a path and returns the relative (using
// ~/path) path and takes an optional boolean to strip the trailing slash.
func MakeRelative(home, dir string, strip bool) string {
	p := strings.Replace(dir, home, "~", 1)

	// should we remove the trailing slash?
	if strip {
		p = strings.TrimRight(p, "/")
	}

	return p
}
