package helpers

import (
	"os"
)

// MkdirIfNotExists is responsible for making a directory
// if the provided directory does not exist. It only
// argument is a dir, which is a path to a dir
func MkdirIfNotExists(dir string) error {
	if dirExists(dir) {
		return nil
	}

	return os.Mkdir(dir, 0755)
}

func dirExists(dir string) bool {
	i, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}

	if i.IsDir() {
		return true
	}

	return false
}
