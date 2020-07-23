package api

import "context"

func (s *NitrodService) NginxService(ctx context.Context, request *NginxServiceRequest) (*ServiceResponse, error) {
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
	_, err := s.command.Run("service", []string{"nginx", action})
	if err != nil {
		s.logger.Println(err)
		return nil, err
	}

	return &ServiceResponse{Message: "successfully " + message + " nginx"}, nil
}
