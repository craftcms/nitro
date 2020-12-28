package api

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/craftcms/nitro/protob"
)

func TestService_Ping(t *testing.T) {
	type fields struct {
		HTTP *http.Client
	}
	type args struct {
		ctx     context.Context
		request *protob.PingRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *protob.PingResponse
		wantErr bool
	}{
		{
			name: "can get a ping",
			args: args{
				ctx:     context.TODO(),
				request: &protob.PingRequest{},
			},
			want:    &protob.PingResponse{Pong: "pong"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				HTTP: tt.fields.HTTP,
			}
			got, err := svc.Ping(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Ping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Ping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Version(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *protob.VersionRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *protob.VersionResponse
		wantErr bool
	}{
		{
			name: "can get the version information from the api",
			args: args{
				ctx:     context.TODO(),
				request: &protob.VersionRequest{},
			},
			want:    &protob.VersionResponse{Version: "develop"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{}
			got, err := svc.Version(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Version() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Version() = %v, want %v", got, tt.want)
			}
		})
	}
}
