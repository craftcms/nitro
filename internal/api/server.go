package api

import (
	"log"
	"net/http"
	"os"
)

type server struct {
	router  *http.ServeMux
	logger  *log.Logger
	service Runner
}

func (s *server) Routes() {
	s.router.HandleFunc("/v1/services/php-fpm", s.handlePhpFpmService())
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// always set the content type
	w.Header().Set("Content-type", "application/json")

	// log the request
	s.logger.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	s.router.ServeHTTP(w, r)
}

func New() *server {
	r := ServiceRunner{}
	return &server{
		router:  http.NewServeMux(),
		logger:  log.New(os.Stdout, "nitrod ", 0),
		service: r,
	}
}
