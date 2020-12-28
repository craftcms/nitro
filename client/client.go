package client

import (
	"fmt"

	"github.com/craftcms/nitro/protob"
	"google.golang.org/grpc"
)

// NewClient is used for generating a new client to interact
// with the gRPC API running in the proxy container
func NewClient(ip, port string) (protob.NitroClient, error) {
	cc, err := grpc.Dial(ip+":"+port, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("unable to create a gRPC client for nitrod, %w", err)
	}

	return protob.NewNitroClient(cc), nil
}
