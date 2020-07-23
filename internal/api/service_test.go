package api

import (
	"context"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

type spyServiceRunner struct {
	Command string
	Args    []string
}

func (r *spyServiceRunner) Run(command string, args []string) ([]byte, error) {
	r.Command = command
	r.Args = args

	return []byte("test"), nil
}

func TestNitrodServer_PhpFpmService(t *testing.T) {
	type fields struct {
		command Runner
		logger  *log.Logger
	}
	type args struct {
		ctx     context.Context
		request *PhpFpmServiceRequest
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        *ServiceResponse
		wantErr     bool
		wantCommand string
		wantArgs    []string
	}{
		{
			name: "can restart php-fpm for version 7.4",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx:     context.TODO(),
				request: &PhpFpmServiceRequest{Version: "7.4", Action: PhpFpmServiceRequest_RESTART},
			},
			want:        &ServiceResponse{Message: "successfully restarted php-fpm 7.4"},
			wantErr:     false,
			wantCommand: "service",
			wantArgs:    []string{"php7.4-fpm", "restart"},
		},
		{
			name: "can stop php-fpm for version 7.3",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx:     context.TODO(),
				request: &PhpFpmServiceRequest{Version: "7.3", Action: PhpFpmServiceRequest_STOP},
			},
			want:        &ServiceResponse{Message: "successfully stopped php-fpm 7.3"},
			wantErr:     false,
			wantCommand: "service",
			wantArgs:    []string{"php7.3-fpm", "stop"},
		},
		{
			name: "can start php-fpm for version 7.2",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx:     context.TODO(),
				request: &PhpFpmServiceRequest{Version: "7.2", Action: PhpFpmServiceRequest_START},
			},
			want:        &ServiceResponse{Message: "successfully started php-fpm 7.2"},
			wantErr:     false,
			wantCommand: "service",
			wantArgs:    []string{"php7.2-fpm", "start"},
		},
		{
			name: "only proper versions for PHP pass validation",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx:     context.TODO(),
				request: &PhpFpmServiceRequest{Version: "7.9", Action: PhpFpmServiceRequest_START},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spy := &spyServiceRunner{}
			s := &NitrodService{
				command: spy,
				logger:  tt.fields.logger,
			}
			got, err := s.PhpFpmService(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("PhpFpmService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PhpFpmService() got = %v, want %v", got, tt.want)
			}

			if tt.wantCommand != "" {
				if spy.Command != tt.wantCommand {
					t.Errorf("wanted the command %v, got %v instead", tt.wantCommand, spy.Command)
				}
			}
			if len(tt.wantArgs) > 0 {
				if !reflect.DeepEqual(spy.Args, tt.wantArgs) {
					t.Errorf("expected the args to be the same\n got %v want %v", spy.Args, tt.wantArgs)
				}
			}
		})
	}
}

func TestNitrodServer_NginxService(t *testing.T) {
	type fields struct {
		command Runner
		logger  *log.Logger
	}
	type args struct {
		ctx     context.Context
		request *NginxServiceRequest
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        *ServiceResponse
		wantErr     bool
		wantCommand string
		wantArgs    []string
	}{
		{
			name: "can restart nginx",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx:     context.TODO(),
				request: &NginxServiceRequest{Action: NginxServiceRequest_RESTART},
			},
			want:        &ServiceResponse{Message: "successfully restarted nginx"},
			wantErr:     false,
			wantCommand: "service",
			wantArgs:    []string{"nginx", "restart"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spy := &spyServiceRunner{}
			s := &NitrodService{
				command: spy,
				logger:  tt.fields.logger,
			}
			got, err := s.NginxService(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("NginxService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NginxService() got = %v, want %v", got, tt.want)
			}
		})
	}
}
