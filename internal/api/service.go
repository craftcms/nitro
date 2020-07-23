package api

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/craftcms/nitro/validate"
)

type NitrodService struct {
	command Runner
	logger  *log.Logger
}

func (s *NitrodService) PhpIniSettings(ctx context.Context, request *ChangePhpIniSettingRequest) (*ServiceResponse, error) {
	// TODO add validation to the value based on the setting
	switch request.GetSetting() {
	case PhpIniSetting_MAX_EXECUTION_TIME:
		_, err := strconv.Atoi(request.GetValue())
		if err != nil {
			return nil, errors.New("max_execution_time must be a valid integer")
		}
	}

	// nitro php iniset max_execution_time val
	// sed -i "/aaa=/c\aaa=xxx" your_file_here
	_, err := s.command.Run("sed", []string{"-i", "s|max_execution_time|max_execution_time = " + request.GetValue() + "|g", "/etc/php/" + request.GetVersion() + "/fpm/php.ini"})
	if err != nil {
		return nil, err
	}
	// todo edit the cli
	// todo restart the php-fpm service

	return &ServiceResponse{Message: "successfully changed the ini setting for max_execution_time to 300"}, nil
}

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

func NewNitrodService() *NitrodService {
	return &NitrodService{
		command: &ServiceRunner{},
		logger:  log.New(os.Stdout, "nitrod ", 0),
	}
}
