package nitrod

import (
	"log"
	"os"
)

var Version string

type Service struct {
	command Runner
	logger  *log.Logger
}

func New() *Service {
	return &Service{
		command: &ServiceRunner{},
		logger:  log.New(os.Stdout, "nitrod ", 0),
	}
}
