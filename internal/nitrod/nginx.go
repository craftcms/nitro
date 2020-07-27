package nitrod

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) NginxService(ctx context.Context, request *NginxServiceRequest) (*ServiceResponse, error) {
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
	if output, err := s.command.Run("service", []string{"nginx", action}); err != nil {
		s.logger.Println(err)
		s.logger.Println("output:", string(output))
		return nil, status.Errorf(codes.Unknown, string(output))
	}

	return &ServiceResponse{Message: "successfully " + message + " nginx"}, nil
}
