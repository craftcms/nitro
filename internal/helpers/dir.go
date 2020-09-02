package helpers

import (
	"os"
)

// MkdirIfNotExists is responsible for making a directory
// if the provided directory does not exist. It only
// argument is a dir, which is a path to a dir
func MkdirIfNotExists(dir string) error {
	if DirExists(dir) {
		return nil
	}

	return os.Mkdir(dir, 0755)
}

// DirExists will return true if the directory exists
func DirExists(dir string) bool {
	i, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}

	if i.IsDir() {
		return true
	}

	return false
}
