package diff

import (
	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

// MountActions takes a machine name and currently attached mounts and mounts from
// the config file. It will then compare the two and determine if we need to add or
// remove mounts from the machine and return the appropriate actions.
func MountActions(name string, attached, file []config.Mount) ([]nitro.Action, error) {
	var actions []nitro.Action

	// if there are more attached than file, we need to run the remove commands
	switch len(file) > (len(attached)) {
	case true:
		for _, mount := range file {
			// check the attached mount if the source already exists
			skip := false

			// check the the source exists in the mount
			for _, mounted := range attached {
				if mounted.AbsSourcePath() == mount.AbsSourcePath() {
					skip = true
				}
			}

			if !skip {
				mountDirAction, err := nitro.MountDir(name, mount.AbsSourcePath(), mount.Dest)
				if err != nil {
					return nil, err
				}
				actions = append(actions, *mountDirAction)
			}
		}

		return actions, nil
	default:
		for _, mount := range attached {
			// check the mount if the source already exists
			skip := false

			// check the the source exists in the mount
			for _, mounted := range file {
				if mounted.Source == mount.Source {
					skip = true
				}
			}

			if !skip {
				mountDirAction, err := nitro.UnmountDir(name, mount.Dest)
				if err != nil {
					return nil, err
				}
				actions = append(actions, *mountDirAction)
			}
		}
	}

	// else run add actions
	return actions, nil
}
