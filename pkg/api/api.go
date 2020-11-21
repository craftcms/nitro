package api

import (
	"context"

	"github.com/craftcms/nitro/pkg/protob"
)

// NewAPI returns an API struct that implements the gRPC API used in the proxy container.
// The gRPC API is used to handle making changes to the Caddy Server via its local API.
func NewAPI() *API {
	return &API{}
}

// API implements the protob.NitroServer interface
type API struct{}

// Ping returns a simple response "pong" from the gRPC API to verify connectivity.
func (a *API) Ping(ctx context.Context, request *protob.PingRequest) (*protob.PingResponse, error) {
	return &protob.PingResponse{Pong: "pong"}, nil
}

// Apply is used to take all of the sites from a Nitro config and apply those changes. The Sites
// in protob.ApplyRequest represents the hostname, aliases (in a comma delimited list), and the
// port for the service. The NGINX container type uses port 8080 and the PHP-FPM container type
// uses port 9000.
func (a *API) Apply(context.Context, *protob.ApplyRequest) (*protob.ApplyResponse, error) {
	return &protob.ApplyResponse{}, nil
}
