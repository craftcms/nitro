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

	// nitro php iniset max_execution_time val
	// sed -i "/aaa=/c\aaa=xxx" your_file_here
	_, err := s.command.Run("sed", []string{"-i", "s|" + setting + "|" + setting + " = " + request.GetValue() + "|g", "/etc/php/" + request.GetVersion() + "/fpm/php.ini"})
	if err != nil {
		return nil, err
	}
	// todo edit the cli
	// todo restart the php-fpm service

	return &ServiceResponse{Message: "successfully changed the ini setting for max_execution_time to 300"}, nil
}
