package nitrod

import (
	"log"
	"os"
)

var Version string

type NitroService struct {
	command Runner
	logger  *log.Logger
}

func NewNitroService() *NitroService {
	return &NitroService{
		command: &ServiceRunner{},
		logger:  log.New(os.Stdout, "nitrod ", 0),
	}
}
