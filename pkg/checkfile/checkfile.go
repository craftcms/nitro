package checkfile

import (
	"fmt"
	"os"
)

var (
	ErrFileNotFound = fmt.Errorf("unable to find the file at the path")
)

// Exists takes a path and verifies if the path or file
// exists. It will return false and an error (ErrFileNotFound)
// if it is unable to stat the path.
func Exists(path string) (bool, error) {
	// make sure the file exists at the path
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, ErrFileNotFound
	}

	return true, nil
}
