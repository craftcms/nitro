package caddyconv

import "github.com/craftcms/nitro/config"

// ToCaddy takes a nitro config struct and converts it to a representation of
// a caddy configuration to send to the Caddy API.
func ToCaddy(sites []config.Site) (*CaddyUpdateRequest, error) {
	// map all of the sites to a server configuration
	req := &CaddyUpdateRequest{}

	// set the defaults for displaying the static page
	req.Srv1.Listen = []string{":80"}

	return req, nil
}
