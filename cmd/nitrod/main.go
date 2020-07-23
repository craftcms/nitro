package main

import (
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/craftcms/nitro/internal/api"
)

func main() {
	port := flag.String("port", "50051", "which port nitro API should listen on")
	flag.Parse()

	lis, err := net.Listen("tcp", "0.0.0.0:"+*port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	api.RegisterNitroServiceServer(s, api.NewNitrodServer())

	if err := s.Serve(lis); err != nil {
		log.Fatal("error when running the server", err)
	}
}
