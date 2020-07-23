package client

import (
	"log"

	"google.golang.org/grpc"

	"github.com/craftcms/nitro/internal/api"
)

// NewClient takes the ip address and port and creates
// a new grpc client for interacting with nitrod api
func NewClient(ip, port string) api.NitroServiceClient {
	// TODO add certificate
	cc, err := grpc.Dial(ip+":"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatal("error creating nitrod client, error:", err)
	}

	return api.NewNitroServiceClient(cc)
}
