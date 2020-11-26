package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/craftcms/nitro/command/version"
	"github.com/craftcms/nitro/protob"
)

const update = `{
  "srv0": {
    "listen": [
      ":443"
    ],
    "routes": [
      {
        "handle": [
          {
            "handler": "subroute",
            "routes": [
              {
                "handle": [
                  {
                    "handler": "reverse_proxy",
                    "upstreams": [
                      {
                        "dial": "updated.nitro:8080"
                      }
                    ]
                  }
                ]
              }
            ]
          }
        ],
        "match": [
          {
            "host": [
              "updated.nitro",
              "updated.test"
            ]
          }
        ],
        "terminal": true
      },
      {
        "handle": [
          {
            "handler": "subroute",
            "routes": [
              {
                "handle": [
                  {
                    "handler": "reverse_proxy",
                    "upstreams": [
                      {
                        "dial": "exampleupdated.nitro:8080"
                      }
                    ]
                  }
                ]
              }
            ]
          }
        ],
        "match": [
          {
            "host": [
              "exampleupdated.nitro",
              "exampleupdated.test"
            ]
          }
        ],
        "terminal": true
      },
      {
        "handle": [
          {
            "handler": "subroute",
            "routes": [
              {
                "handle": [
                  {
                    "handler": "vars",
                    "root": "/var/www/html"
                  },
                  {
                    "handler": "file_server",
                    "hide": [
                      "/etc/caddy/Caddyfile"
                    ]
                  }
                ]
              }
            ]
          }
        ],
        "terminal": true
      }
    ]
  },
  "srv1": {
    "listen": [
      ":80"
    ],
    "routes": [
      {
        "handle": [
          {
            "handler": "vars",
            "root": "/var/www/html"
          },
          {
            "handler": "file_server",
            "hide": [
              "/etc/caddy/Caddyfile"
            ]
          }
        ]
      }
    ]
  }
}`

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
func (a *API) Apply(ctx context.Context, request *protob.ApplyRequest) (*protob.ApplyResponse, error) {
	for k, v := range request.GetSites() {
		fmt.Println(k, v)
	}
	resp := &protob.ApplyResponse{}

	// TODO(jasonmccallister) create the struct to represent the servers to send
	// to the caddy api.

	// send the update
	res, err := http.Post("http://127.0.0.1:2019/config/apps/http/servers", "application/json", strings.NewReader(update))
	if err != nil {
		resp.Message = "error updating the Caddy API"
		resp.Error = true

		return resp, err
	}

	// check the status code
	if res.StatusCode != http.StatusOK {
		resp.Message = fmt.Sprintf("received %d response from caddy api", res.StatusCode)
		resp.Error = true

		return resp, nil
	}

	resp.Message = "success"
	resp.Error = false

	return resp, nil
}

// Version is used to check the container image version with the CLI version
func (a *API) Version(context.Context, *protob.VersionRequest) (*protob.VersionResponse, error) {
	return &protob.VersionResponse{Version: version.Version}, nil
}
