package client

import (
	"context"
	"fmt"
)

func Composer(ctx context.Context, absDir, version, action string) error {
	// TODO check if there is a composer.json

	// TODO create the temp container

	// get the version from the flag, default to 1
	switch action {
	case "update":
		fmt.Print("Running composer update with version", version)
		// docker run --rm --interactive --tty --volume absDir:/app composer:version update --ignore-platform-reqs
	default:
		fmt.Print("Running composer install with version", version)
		// docker run --rm --interactive --tty --volume absDir:/app composer:version install --ignore-platform-reqs
	}

	return nil
}
