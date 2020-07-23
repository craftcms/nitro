package api

import (
	"context"
	"log"
	"os"

	"github.com/craftcms/nitro/validate"
)

type NitrodServer struct {
	command Runner
	logger  *log.Logger
}

func (s *NitrodServer) PhpFpmService(ctx context.Context, request *PhpFpmServiceRequest) (*ServiceResponse, error) {
	// validate the request
	if err := validate.PHPVersion(request.GetVersion()); err != nil {
		s.logger.Println(err)
		return nil, err
	}

	var action string
	var message string
	switch request.GetAction() {
	case PhpFpmServiceRequest_START:
		message = "started"
		action = "start"
	case PhpFpmServiceRequest_STOP:
		message = "stopped"
		action = "stop"
	default:
		message = "restarted"
		action = "restart"
	}

	// perform the action on the php-fpm service
	_, err := s.command.Run("service", []string{"php" + request.GetVersion() + "-fpm", action})
	if err != nil {
		s.logger.Println(err)
		return nil, err
	}

	return &ServiceResponse{Message: "successfully " + message + " php-fpm " + request.GetVersion()}, nil
}

func (s *NitrodServer) NginxService(ctx context.Context, request *NginxServiceRequest) (*ServiceResponse, error) {
	var action string
	var message string
	switch request.GetAction() {
	case NginxServiceRequest_START:
		message = "started"
		action = "start"
	case NginxServiceRequest_STOP:
		message = "stopped"
		action = "stop"
	default:
		message = "restarted"
		action = "restart"
	}

	// perform the action on the php-fpm service
	_, err := s.command.Run("service", []string{"nginx", action})
	if err != nil {
		s.logger.Println(err)
		return nil, err
	}

	return &ServiceResponse{Message: "successfully " + message + " nginx"}, nil
}

func NewNitrodServer() *NitrodServer {
	return &NitrodServer{
		command: &ServiceRunner{},
		logger:  log.New(os.Stdout, "nitrod ", 0),
	}
}
