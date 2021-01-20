package checkfile

import (
	"context"
	"fmt"
	"os"
)

var (
	ErrFileNotFound = fmt.Errorf("unable to find the file at the path")
)

func Exists(ctx context.Context, path string) (bool, error) {
	// make sure the file exists at the path
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, ErrFileNotFound
	}

	return true, nil
}
