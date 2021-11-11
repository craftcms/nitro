package bindmounts

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	config "github.com/craftcms/nitro/pkg/config/v3"
)

func ForApp(app config.App) ([]string, error) {
	if len(app.Excludes) > 0 {

	}

	return nil, fmt.Errorf("we are not done")
}

// FromDir takes a directory path and a list of directories to exclude and returns the
// absolute path to each directory and file that should be bind mounted into a container.
func FromDir(path string, excludes []string) []string {
	var mounts []string

	//read the directory
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}

	// loop over each file/dir in provided directory
	for _, d := range dirs {
		var excluded bool
		for _, e := range excludes {
			if d.Name() == e {
				excluded = true
				break
			}
		}

		if !excluded {
			mounts = append(mounts, filepath.Join(path, d.Name()))
		}
	}

	return mounts
}
