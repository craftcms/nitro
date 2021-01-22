package pathexists

import (
	"os"
)

// IsDirectory takes a path and returns true if the
// provided path is a directory.
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

// IsFile takes a path an verifies the path exists
// and is a file.
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.Mode().IsRegular()
}
