package nitrod

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/craftcms/nitro/internal/validate"
)

func (s *NitroService) GetPhpIniSetting(ctx context.Context, req *GetPhpIniSettingRequest) (*ServiceResponse, error) {
	// validate php version
	if err := validate.PHPVersion(req.GetVersion()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// get the setting from the php ini_get function
	output, err := s.command.Run("php"+req.GetVersion(), []string{"-r", fmt.Sprintf("echo ini_get('%s');", req.GetSetting())})
	if err != nil {
		s.logger.Println("error getting ini setting:", err)
		s.logger.Println("output:", string(output))
		return nil, status.Errorf(codes.Unknown, string(output))
	}

	return &ServiceResponse{Message: fmt.Sprintf("The setting %q is currently set to %v", req.GetSetting(), string(output))}, nil
}
