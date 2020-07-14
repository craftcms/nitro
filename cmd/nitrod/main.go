package main

import (
	"flag"
	"github.com/craftcms/nitro/internal/api"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "9999", "which port the nitro API should listen on")

	flag.Parse()

	srv := api.New()

	srv.Routes()

	log.Println("listening on port", *port)

	log.Fatal(http.ListenAndServe(":"+*port, srv))
}
