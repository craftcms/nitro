package nitrod

import (
	"context"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

func TestNitroService_DisableXdebug(t *testing.T) {
	type fields struct {
		command Runner
		logger  *log.Logger
	}
	type args struct {
		ctx context.Context
		req *DisableXdebugRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		want         *ServiceResponse
		wantErr      bool
		wantCommands []string
		wantArgs     []map[string][]string
	}{
		{
			name: "can disable xdebug for php version 7.4",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx: context.TODO(),
				req: &DisableXdebugRequest{Version: "7.4"},
			},
			want:         &ServiceResponse{Message: "Disabled xdebug for PHP 7.4"},
			wantErr:      false,
			wantCommands: []string{"phpdismod", "service"},
			wantArgs: []map[string][]string{
				{
					"phpdismod": {"-v", "7.4", "xdebug"},
				},
				{
					"service": {"php7.4-fpm", "restart"},
				},
			},
		},
		{
			name: "request fails validation",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx: context.TODO(),
				req: &DisableXdebugRequest{Version: "7777"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spy := &spyChainRunner{}
			s := &NitroService{
				command: spy,
				logger:  tt.fields.logger,
			}
			got, err := s.DisableXdebug(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DisableXdebug() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DisableXdebug() got = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(spy.Commands, tt.wantCommands) {
				t.Errorf("expected the commands to be:\n%v\n, got:\n%v", tt.wantCommands, spy.Commands)
			}

			if !reflect.DeepEqual(spy.Args, tt.wantArgs) {
				t.Errorf("expected the args to be:\n%v\ngot:\n%v", tt.wantArgs, spy.Args)
			}
		})
	}
}

func TestNitroService_EnableXdebug(t *testing.T) {
	type fields struct {
		command Runner
		logger  *log.Logger
	}
	type args struct {
		ctx context.Context
		req *EnableXdebugRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		want         *ServiceResponse
		wantErr      bool
		wantCommands []string
		wantArgs     []map[string][]string
	}{
		{
			name: "can enable xdebug for php version 7.4",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx: context.TODO(),
				req: &EnableXdebugRequest{Version: "7.4"},
			},
			want:         &ServiceResponse{Message: "Enabled xdebug for PHP 7.4"},
			wantErr:      false,
			wantCommands: []string{"phpenmod", "service"},
			wantArgs: []map[string][]string{
				{
					"phpenmod": {"-v", "7.4", "xdebug"},
				},
				{
					"service": {"php7.4-fpm", "restart"},
				},
			},
		},
		{
			name: "request fails validation",
			fields: fields{
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			args: args{
				ctx: context.TODO(),
				req: &EnableXdebugRequest{Version: "7777"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spy := &spyChainRunner{}
			s := &NitroService{
				command: spy,
				logger:  tt.fields.logger,
			}
			got, err := s.EnableXdebug(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnableXdebug() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnableXdebug() got = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(spy.Commands, tt.wantCommands) {
				t.Errorf("expected the commands to be:\n%v\n, got:\n%v", tt.wantCommands, spy.Commands)
			}

			if !reflect.DeepEqual(spy.Args, tt.wantArgs) {
				t.Errorf("expected the args to be:\n%v\ngot:\n%v", tt.wantArgs, spy.Args)
			}
		})
	}
}


