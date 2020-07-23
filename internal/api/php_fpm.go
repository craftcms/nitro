package api

import (
	"context"

	"github.com/craftcms/nitro/validate"
)

func (s *NitrodService) PhpFpmService(ctx context.Context, request *PhpFpmServiceRequest) (*ServiceResponse, error) {
	// validate the request
	if err := validate.PHPVersion(request.GetVersion()); err != nil {
		s.logger.Println(err)
		return nil, err
	}

	var action string
	var message string
	switch request.GetAction() {
	case ServiceAction_START:
		message = "started"
		action = "start"
	case ServiceAction_STOP:
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
