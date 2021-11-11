package bindmounts

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/paths"
)

func ForApp(app config.App, home string) ([]string, error) {
	path, err := paths.Clean(home, app.Path)
	if err != nil {
		return nil, err
	}

	if len(app.Excludes) > 0 {
		var binds []string
		for _, v := range FromDir(path, app.Excludes) {
			_, f := filepath.Split(v)

			binds = append(binds, fmt.Sprintf("%s:/app/%s:rw", v, f))
		}

		return binds, nil
	}

	// return the entire directory as the bind mount
	return []string{fmt.Sprintf("%s:/app:rw", path)}, nil
}

func ForSite(s config.Site, home string) ([]string, error) {
	// get the abs path for the site
	path, err := s.GetAbsPath(home)
	if err != nil {
		return nil, err
	}

	// are there files or directories we should exclude?
	if len(s.Excludes) > 0 {
		var binds []string
		for _, v := range FromDir(path, s.Excludes) {
			_, f := filepath.Split(v)

			binds = append(binds, fmt.Sprintf("%s:/app/%s:rw", v, f))
		}

		return binds, nil
	}

	// return the entire directory as the bind mount
	return []string{fmt.Sprintf("%s:/app:rw", path)}, nil
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
