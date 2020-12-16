package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/craftcms/nitro/command/version"
	"github.com/craftcms/nitro/pkg/caddyconv"
	"github.com/craftcms/nitro/protob"
)

// NewAPI returns an API struct that implements the gRPC API used in the proxy container.
// The gRPC API is used to handle making changes to the Caddy Server via its local API.
func NewAPI() *API {
	return &API{
		Client: http.DefaultClient,
	}
}

// API implements the protob.NitroServer interface
type API struct {
	Client *http.Client
}

// Ping returns a simple response "pong" from the gRPC API to verify connectivity.
func (a *API) Ping(ctx context.Context, request *protob.PingRequest) (*protob.PingResponse, error) {
	return &protob.PingResponse{Pong: "pong"}, nil
}

// Apply is used to take all of the sites from a Nitro config and apply those changes. The Sites
// in protob.ApplyRequest represents the hostname, aliases (in a comma delimited list), and the
// port for the service. The NGINX container type uses port 8080 and the PHP-FPM container type
// uses port 9000.
func (a *API) Apply(ctx context.Context, request *protob.ApplyRequest) (*protob.ApplyResponse, error) {
	resp := &protob.ApplyResponse{}

	// if there is no client, use the default
	if a.Client == nil {
		a.Client = http.DefaultClient
	}

	// convert each of the sites into a route
	routes := []caddyconv.ServerRoute{}
	for k, site := range request.GetSites() {
		// get all of the host names for the site
		hosts := []string{site.GetHostname()}
		if site.GetAliases() != "" {
			hosts = append(hosts, strings.Split(site.GetAliases(), ",")...)
		}

		// create the route for each of the sites
		routes = append(routes, caddyconv.ServerRoute{
			Handle: []caddyconv.RouteHandle{
				{
					Handler: "reverse_proxy",
					Upstreams: []caddyconv.Upstream{
						{
							Dial: fmt.Sprintf("%s:%d", k, site.GetPort()),
						},
					},
				},
			},
			Match: []caddyconv.Match{
				{
					Host: hosts,
				},
			},
			Terminal: true,
		})
	}

	update := caddyconv.CaddyUpdateRequest{}

	// add the routes to the first server
	update.Srv0 = caddyconv.Server{
		Listen: []string{":443", ":9003"},
		Routes: routes,
	}

	// set the default welcome server
	update.Srv1 = caddyconv.Server{
		Listen: []string{":80"},
		Routes: []caddyconv.ServerRoute{
			{
				Handle: []caddyconv.RouteHandle{
					{
						Handler: "vars",
						Root:    "/var/www/html",
					},
					{
						Handler: "file_server",
						Root:    "/var/www/html",
						Hide:    []string{"/etc/caddy/Caddyfile"},
					},
				},
				Terminal: true,
			},
		},
	}

	content, err := json.Marshal(&update)
	if err != nil {
		return nil, err
	}

	// send the update
	res, err := a.Client.Post("http://127.0.0.1:2019/config/apps/http/servers", "application/json", bytes.NewReader(content))
	if err != nil {
		resp.Message = "error updating Caddy API"
		resp.Error = true

		return resp, err
	}

	// check the status code
	if res.StatusCode != http.StatusOK {
		resp.Message = fmt.Sprintf("received %d response from caddy api", res.StatusCode)
		resp.Error = true

		return resp, nil
	}

	// set the message and error to false
	resp.Message = fmt.Sprintf("successfully applied changes, sites: %d", len(request.GetSites()))
	resp.Error = false

	return resp, nil
}

// Version is used to check the container image version with the CLI version
func (a *API) Version(context.Context, *protob.VersionRequest) (*protob.VersionResponse, error) {
	return &protob.VersionResponse{Version: version.Version}, nil
}
