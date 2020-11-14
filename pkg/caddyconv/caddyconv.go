package caddyconv

import "github.com/craftcms/nitro/internal/config"

// ToCaddy takes a nitro config struct and converts it to a representation of
// a caddy configuration to send to the Caddy API.
func ToCaddy(config *config.Config) (*CaddyConfig, error) {
	// map all of the sites to a server configuration
	servers := map[string]CaddyServer{}
	for _, site := range config.Sites {
		servers[site.Hostname] = CaddyServer{
			// TODO(jasonmccallister) make this grab the sites aliases
			Listen: []string{":443"},
			// Routes:
		}
	}

	// set the defaults
	caddy := &CaddyConfig{
		Admin: CaddyAdmin{
			Listen: "localhost:2019",
		},
	}
	// set the tls configuration
	// tls := CaddyAppsTLS{}
	// generate

	return caddy, nil
}
