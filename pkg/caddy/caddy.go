package caddy

type UpdateRequest struct {
	Srv0 Server `json:"srv0,omitempty"`
	Srv1 Server `json:"srv1,omitempty"`
}

type Server struct {
	Listen []string      `json:"listen"`
	Routes []ServerRoute `json:"routes"`
}

type ServerRoute struct {
	Handle   []RouteHandle `json:"handle"`
	Match    []Match       `json:"match,omitempty"`
	Terminal bool          `json:"terminal"`
}

type RouteHandle struct {
	Handler   string     `json:"handler"`
	Root      string     `json:"root,omitempty"`
	Upstreams []Upstream `json:"upstreams,omitempty"`
	Hide      []string   `json:"hide,omitempty"`
}

type Match struct {
	Host []string `json:"host"`
}

type Upstream struct {
	Dial string `json:"dial,omitempty"`
}
