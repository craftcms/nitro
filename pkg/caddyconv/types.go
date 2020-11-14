package caddyconv

type CaddyConfig struct {
	Admin CaddyAdmin `json:"admin"`
	Apps  struct {
		HTTP struct {
			Servers map[string]CaddyServer `json:"servers"`
		} `json:"http"`
		TLS CaddyAppsTLS `json:"tls"`
	} `json:"apps"`
}

type CaddyAppsTLS struct {
	Automation struct {
		Policies []struct {
			Issuer struct {
				Module string `json:"module"`
			} `json:"issuer"`
		} `json:"policies"`
	} `json:"automation"`
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

type CaddyRoute struct {
}

type CaddyAdmin struct {
	Listen string `json:"listen"`
}
