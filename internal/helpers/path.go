package helpers

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// PathName takes a path and will return the directory name of the
// parent dir (e.g. /nitro/sites/test will return test)
func PathName(path string) (string, error) {
	p, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	s := strings.Split(p, string(os.PathSeparator))

	if len(s) < 2 {
		return "", errors.New("unexpected wrong number of paths")
	}

	return s[len(s)-1], nil
}
