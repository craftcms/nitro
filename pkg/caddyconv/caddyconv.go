package caddyconv

import "github.com/craftcms/nitro/internal/config"

type CaddyConfig struct {
	Admin struct {
		Listen string `json:"listen"`
	} `json:"admin"`
	Apps struct {
		HTTP struct {
			Servers map[string][]CaddyServer `json:"servers"`
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
				Group  string `json:"group,omitempty"`
				Handle []struct {
					Handler string `json:"handler"`
					URI     string `json:"uri"`
				} `json:"handle"`
				Match []struct {
					File struct {
						TryFiles []string `json:"try_files"`
					} `json:"file"`
					Not []struct {
						Path []string `json:"path"`
					} `json:"not"`
				} `json:"match,omitempty"`
			} `json:"routes"`
		} `json:"handle"`
		Terminal bool `json:"terminal"`
	} `json:"routes"`
}

func ToJSON(config config.Config) {

}
