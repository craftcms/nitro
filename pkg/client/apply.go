package client

import (
	"context"
	"fmt"
)

// Apply is used to create a
func (cli *Client) Apply(ctx context.Context, php string) error {
	// TODO(jasonmccallister) get all of the sites, their local path, the php version, and the type of project (nginx or PHP-FPM)
	// for range in sites
	// does the container exist and the php version/type match?
	// else if not recreate

	//TODO(jasonmccallister) get all of the databases, engine, version, and ports and create a container for each

	// TODO(jasonmccallister) convert the sites into a Caddy json config and send to the API

	return fmt.Errorf("not yet implemented")
}
