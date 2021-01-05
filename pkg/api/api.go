package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/craftcms/nitro/command/version"
	"github.com/craftcms/nitro/pkg/caddy"
	"github.com/craftcms/nitro/protob"
)

// NewService takes the address to the Caddy API and returns an API struct that
// implements the gRPC API used in the proxy container. The gRPC API is used to
// handle making changes to the Caddy Server via its local API. If no addr is
// provided, it will set the default addr to http://127.0.0.1:2019
func NewService(addr string) *Service {
	return &Service{
		Addr: addr,
		HTTP: http.DefaultClient,
	}
}

// Service implements the protob.NitroServer interface
type Service struct {
	Addr string
	HTTP *http.Client
}

// Ping returns a simple response "pong" from the gRPC API to verify connectivity.
func (svc *Service) Ping(ctx context.Context, request *protob.PingRequest) (*protob.PingResponse, error) {
	return &protob.PingResponse{Pong: "pong"}, nil
}

// Apply is used to take all of the sites from a Nitro config and apply those changes. The Sites
// in protob.ApplyRequest represents the hostname, aliases (in a comma delimited list), and the
// port for the service. The NGINX container type uses port 8080 and the PHP-FPM container type
// uses port 9000.
func (svc *Service) Apply(ctx context.Context, request *protob.ApplyRequest) (*protob.ApplyResponse, error) {
	// if there is no client, use the default
	if svc.HTTP == nil {
		svc.HTTP = http.DefaultClient
	}

	// set the addr if not provided
	if svc.Addr == "" {
		svc.Addr = "http://127.0.0.1:2019"
	}

	// convert each of the sites into a route
	routes := []caddy.ServerRoute{}
	for k, site := range request.GetSites() {
		// get all of the host names for the site
		hosts := []string{site.GetHostname()}
		if site.GetAliases() != "" {
			hosts = append(hosts, strings.Split(site.GetAliases(), ",")...)
		}

		// create the route for each of the sites
		routes = append(routes, caddy.ServerRoute{
			Handle: []caddy.RouteHandle{
				{
					Handler: "reverse_proxy",
					Upstreams: []caddy.Upstream{
						{
							Dial: fmt.Sprintf("%s:%d", k, site.GetPort()),
						},
					},
				},
			},
			Match: []caddy.Match{
				{
					Host: hosts,
				},
			},
			Terminal: true,
		})
	}

	update := caddy.UpdateRequest{}

	// add the routes to the first server
	update.Srv0 = caddy.Server{
		Listen: []string{":443"},
		Routes: routes,
	}

	// set the default welcome server
	update.Srv1 = caddy.Server{
		Listen: []string{":80"},
		Routes: append(routes, caddy.ServerRoute{
			Handle: []caddy.RouteHandle{
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
		}),
	}

	content, err := json.Marshal(&update)
	if err != nil {
		return nil, err
	}

	// send the update
	res, err := svc.HTTP.Post(svc.Addr+"/config/apps/http/servers", "application/json", bytes.NewReader(content))
	if err != nil {
		return &protob.ApplyResponse{
			Message: fmt.Sprintf("Error updating Caddy API, err: %s", err.Error()),
			Error:   true,
		}, err
	}

	// check the status code
	if res.StatusCode != http.StatusOK {
		return &protob.ApplyResponse{
			Message: fmt.Sprintf("Received %d response from Caddy API", res.StatusCode),
			Error:   true,
		}, nil
	}

	// set the message and error to false
	return &protob.ApplyResponse{
		Message: fmt.Sprintf("Successfully applied changes, sites: %d", len(request.GetSites())),
		Error:   false,
	}, nil
}

// Version is used to check the container image version with the CLI version
func (svc *Service) Version(ctx context.Context, request *protob.VersionRequest) (*protob.VersionResponse, error) {
	return &protob.VersionResponse{Version: version.Version}, nil
}
