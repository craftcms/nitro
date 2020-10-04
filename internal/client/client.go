package client

import (
	"log"

	"google.golang.org/grpc"

	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/nitrod"
)

// NewClient takes the ip address and port and creates
// a new grpc client for interacting with nitrod nitrod
// service.
func NewClient(ip, port string) (nitrod.NitroServiceClient, error) {
	cc, err := grpc.Dial(ip+":"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatal("error creating nitrod client, error:", err)
	}

	return nitrod.NewNitroServiceClient(cc), nil
}

func NewDefaultClient(machine string) (nitrod.NitroServiceClient, error) {
	ip := nitro.IP(machine, nitro.NewMultipassRunner("multipass"))

	cc, err := grpc.Dial(ip+":"+"50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("error creating nitrod client, error:", err)
	}

	return nitrod.NewNitroServiceClient(cc), nil
}

// NewSystemClient takes the ip address and port and creates
// a new gRPC client for interacting with the nitrod systems
// service.
func NewSystemClient(ip, port string) (nitrod.SystemServiceClient, error) {
	cc, err := grpc.Dial(ip+":"+port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return nitrod.NewSystemServiceClient(cc), nil
}
