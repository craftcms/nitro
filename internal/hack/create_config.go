package hack

import (
	"github.com/craftcms/nitro/internal/helpers"
)

func CreateConfig(home, machine string) error {
	nitroDir := home + "/.nitro/"

	if err := helpers.MkdirIfNotExists(nitroDir); err != nil {
		return err
	}

	return helpers.CreateFileIfNotExist(nitroDir + machine + ".yaml")
}
