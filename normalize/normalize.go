package normalize

import (
	"os"
	"path/filepath"
	"strings"
)

// Path is responsible for taking a path, relative
// or otherwise, and returning the name of the file,
// the absolute path, and an error if not found.
func Path(path, home string) (string, string, error) {
	p := strings.Split(path, string(os.PathSeparator))

	if strings.Contains(p[0], "~") {p[0] = home
		path = strings.Join(p, string(os.PathSeparator))
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", "", err
	}

	filename := p[len(p)-1]

	return filename, absPath, nil
}
