package nitrod

import (
	"context"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

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
				request: &NginxServiceRequest{Action: ServiceAction_RESTART},
			},
			want:        &ServiceResponse{Message: "Successfully restarted nginx"},
			wantErr:     false,
			wantCommand: "service",
			wantArgs:    []string{"nginx", "restart"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spy := &spyServiceRunner{}
			s := &SystemService{
				command: spy,
				logger:  tt.fields.logger,
			}
			got, err := s.Nginx(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("NginxService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NginxService() got = %v, want %v", got, tt.want)
			}

			if tt.wantCommand != "" {
				if spy.Command != tt.wantCommand {
					t.Errorf("wanted the command %v, got %v instead", tt.wantCommand, spy.Command)
				}
			}
			if len(tt.wantArgs) > 0 {
				if !reflect.DeepEqual(spy.Args, tt.wantArgs) {
					t.Errorf("expected the args to be the same\n got:\n%v\n want:\n%v", spy.Args, tt.wantArgs)
				}
			}
		})
	}
}
