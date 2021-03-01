package directory

import (
	"io"
	"os"
)

// IsEmpty takes a directory and will verify the if the directory
// is empty or not.
func IsEmpty(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1)

	return err == io.EOF
}
