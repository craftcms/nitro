package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/craftcms/nitro/internal/api"
)

func main() {
	port := flag.String("port", "9999", "which port the nitro API should listen on")
	flag.Parse()

	srv := api.New()

	srv.Routes()

	log.Println("listening on port", *port)

	log.Fatal(http.ListenAndServe(":"+*port, srv))
}
