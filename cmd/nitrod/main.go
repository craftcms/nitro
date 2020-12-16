package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/craftcms/nitro/api"
	"github.com/craftcms/nitro/protob"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetPrefix("nitrod")

	// assign the port as a flag with a default
	port := flag.String("port", "5000", "which port API should listen on")
	flag.Parse()

	// create the network listener
	lis, err := net.Listen("tcp", "0.0.0.0:"+*port)
	if err != nil {
		log.Fatal(err)
	}

	// create the grpc server
	s := grpc.NewServer()

	protob.RegisterNitroServer(s, api.NewAPI())

	fmt.Println("gRPC API listening on port", *port)

	// server the grpc service
	if err := s.Serve(lis); err != nil {
		log.Fatal("error when running the api", err)
	}
}
