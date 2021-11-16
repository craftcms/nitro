package appaware

import (
	"fmt"
	"strings"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/paths"
)

// Detect takes a config and a working directory and returns the hostname
// of the app if we can detect it otherwise it will return an error.
func Detect(cfg config.Config, dir string) (string, error) {
	// get the current path of the command
	current, err := paths.Clean(cfg.HomeDirectory, dir)
	if err != nil {
		return "", err
	}

	// check the apps and try to find the matching app
	for _, app := range cfg.ParsedApps {
		path, err := paths.Clean(cfg.HomeDirectory, app.Path)
		if err != nil {
			return "", err
		}

		// does the path match?
		if strings.Contains(current, path) {
			return app.Hostname, nil
		}
	}

	return "", fmt.Errorf("no app detected at %s", current)
}
