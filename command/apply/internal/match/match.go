package match

import (
	"fmt"
	"strings"

	"github.com/craftcms/nitro/config"
	"github.com/docker/docker/api/types"
)

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
func Site(home string, site config.Site, php config.PHP, container types.ContainerJSON) bool {
	// check if the image does not match
	if fmt.Sprintf("docker.io/craftcms/nginx:%s-dev", site.PHP) != container.Config.Image {
		return false
	}

	// check the environment variables
	for _, e := range container.Config.Env {
		sp := strings.Split(e, "=")

		// show only the environment variables we know about/support
		if _, ok := config.Envs[sp[0]]; ok {
			// TODO(jasonmccallister) check if the value matches the config option
			fmt.Println(sp[0])
		}
	}

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
