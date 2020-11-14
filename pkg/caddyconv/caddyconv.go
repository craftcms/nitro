package caddyconv

import "github.com/craftcms/nitro/internal/config"

type CaddyConfig struct {
	Admin CaddyAdmin `json:"admin"`
	Apps  struct {
		HTTP struct {
			Servers map[string]CaddyServer `json:"servers"`
		} `json:"http"`
		TLS struct {
			Automation struct {
				Policies []struct {
					Issuer struct {
						Module string `json:"module"`
					} `json:"issuer"`
				} `json:"policies"`
			} `json:"automation"`
		} `json:"tls"`
	} `json:"apps"`
}

type CaddyServer struct {
	Listen []string `json:"listen"`
	Routes []struct {
		Match []struct {
			Host []string `json:"host"`
		} `json:"match,omitempty"`
		Handle []struct {
			Handler string `json:"handler"`
			Routes  []struct {
				Handle []struct {
					Handler   string `json:"handler"`
					Upstreams []struct {
						Dial string `json:"dial"`
					} `json:"upstreams"`
				} `json:"handle"`
			} `json:"routes"`
		} `json:"handle"`
		Terminal bool `json:"terminal"`
	} `json:"routes"`
}

type CaddyAdmin struct {
	Listen string `json:"listen"`
}

func ToCaddy(config *config.Config) (*CaddyConfig, error) {
	caddy := &CaddyConfig{Admin: CaddyAdmin{Listen: "localhost:2019"}}
	// generate

	return caddy, nil
}
