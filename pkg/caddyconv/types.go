package caddyconv

type CaddyUpdateRequest struct {
	Srv0 struct {
		Listen []string `json:"listen"`
		Routes []struct {
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
			Match []struct {
				Host []string `json:"host"`
			} `json:"match,omitempty"`
			Terminal bool `json:"terminal"`
		} `json:"routes"`
	} `json:"srv0"`
	Srv1 struct {
		Listen []string `json:"listen"`
		Routes []Route `json:"routes"`
	} `json:"srv1"`
}

type Route struct {
	Handle []Handler `json:"handle"`
}

type Handler struct {
	Handler string   `json:"handler"`
	Root    string   `json:"root,omitempty"`
	Hide    []string `json:"hide,omitempty"`
}
