package helpers

import (
	"errors"
	"os"
)

func CreateFileIfNotExist(filename string) error {
	// does it exist
	if FileExists(filename) {
		return errors.New("file already exists")
	}

	// create it
	_, err := os.Create(filename)

	return err
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
