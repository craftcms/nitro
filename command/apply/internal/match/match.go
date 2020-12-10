package match

import (
	"fmt"

	"github.com/craftcms/nitro/config"
	"github.com/docker/docker/api/types"
)

// TODO(jasonmccallister) determine if we need this
type containerInstance struct {
	Hostname  string
	Image     string
	Container types.Container
	Envs      []string
}

// Mounts takes a containers existing mounts against the sites expected mounts
// and returns true if the mounts do not match.
func Mounts(existing []types.MountPoint, expected map[string]string) bool {
	if len(existing) != len(expected) {
		return false
	}

	return true
}

// Site takes the home directory, site, php config, and a container to determine if they
// match whats expected.
func Site(home string, site config.Site, php config.PHP, container types.Container) bool {
	// check if the image does not match
	if fmt.Sprintf("docker.io/craftcms/nginx:%s-dev", site.PHP) != container.Image {
		return false
	}

	// check the environment variables

	// check the mounts

	// get the main site path (e.g. ~/dev/craft-dev)
	path, err := site.GetAbsPath(home)
	if err != nil {
		return false
	}

	// get any additional mounts for the site (e.g. mounts:)
	expected, err := site.GetAbsMountPaths(home)
	if err != nil {
		return false
	}

	// hard code the path to the first site mount (which is the path)
	expected[path] = "/app"

	// check if the mounts exist
	if Mounts(container.Mounts, expected) == false {
		return false
	}

	return true
}
