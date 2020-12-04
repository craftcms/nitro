package match

import (
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
