package helpers

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// PathName takes a path and will return the directory name of the
// dir (e.g. /nitro/sites/test will return test)
func PathName(path string) (string, error) {
	p, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	s := strings.Split(p, "/")

	if len(s) < 2 {
		return "", errors.New("unexpected wrong number of paths")
	}

	return s[len(s)-1], nil
}

// GetDirectoryArg takes an argument from the command
// and will return the name of the directory with
// the full path to the directory. Mostly used
// on the `nitro add` command.
func GetDirectoryArg(args []string) (string, string, error) {
	// always get the current directory
	wd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	// if there was an argument provided get that working directory
	if len(args) > 0 {
		// check if this is a relative pth
		switch strings.HasPrefix(args[0], "~") {
		case true:
			home, err := homedir.Dir()
			if err != nil {
				return "", "", err
			}

			wd = strings.Replace(args[0], "~", home, 1)
		default:
			// get the abs path
			wd, err = filepath.Abs(args[0])
			if err != nil {
				return "", "", err
			}
		}
	}

	path := strings.Split(wd, string(os.PathSeparator))

	return RemoveTrailingSlash(path[len(path)-1]), wd, nil
}
