package api

import (
	"context"
	"errors"
	"strconv"
)

func (s *NitrodService) PhpIniSettings(ctx context.Context, request *ChangePhpIniSettingRequest) (*ServiceResponse, error) {
	// TODO add validation to the value based on the setting
	var setting string
	switch request.GetSetting() {
	case PhpIniSetting_MAX_EXECUTION_TIME:
		_, err := strconv.Atoi(request.GetValue())
		if err != nil {
			return nil, errors.New("max_execution_time must be a valid integer")
		}
		setting = "max_execution_time"
	default:
		e := errors.New("changing this setting is not authorized")
		s.logger.Println(e)
		return nil, e
	}

	_, err := s.command.Run("sed", []string{"-i", "s|" + setting + "|" + setting + " = " + request.GetValue() + "|g", "/etc/php/" + request.GetVersion() + "/fpm/php.ini"})
	if err != nil {
		s.logger.Println("error changing ini setting, error:", err)
		return nil, err
	}

	_, err = s.command.Run("service", []string{"php" + request.GetVersion() + "-fpm", "restart"})
	if err != nil {
		s.logger.Println("error restarting php-fpm, error:", err)
		return nil, err
	}

	return &ServiceResponse{Message: "successfully changed the ini setting for " + setting + " to " + request.GetValue()}, nil
}
