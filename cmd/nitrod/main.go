package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/craftcms/nitro/internal/nitrod"
)

func main() {
	// assign the port as a flag with a default
	port := flag.String("port", "50051", "which port nitro API should listen on")
	flag.Parse()

	// create the network listener
	lis, err := net.Listen("tcp", "0.0.0.0:"+*port)
	if err != nil {
		log.Fatal(err)
	}

	// create the grpc server
	s := grpc.NewServer()

	// register our services
	nitrod.RegisterNitroServiceServer(s, nitrod.NewNitroService())
	nitrod.RegisterSystemServiceServer(s, nitrod.NewSystemService())

	fmt.Println("running nitrod on port", *port)

	// server the grpc service
	if err := s.Serve(lis); err != nil {
		log.Fatal("error when running the server", err)
	}
}
