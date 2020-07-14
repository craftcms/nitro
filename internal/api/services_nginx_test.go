package api

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func Test_server_handleNginxService(t *testing.T) {
	type fields struct {
		router *http.ServeMux
		logger *log.Logger
	}
	tests := []struct {
		name        string
		fields      fields
		body        io.Reader
		statusCode  int
		wantCommand string
		wantArgs    []string
	}{
		{
			name: "testing request validation",
			fields: fields{
				router: http.NewServeMux(),
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "testing valid requests call correct actions",
			fields: fields{
				router: http.NewServeMux(),
				logger: log.New(ioutil.Discard, "testing", 0),
			},
			body:        strings.NewReader(`{"action":"restart"}`),
			wantCommand: "service",
			wantArgs:    []string{"nginx", "restart"},
			statusCode:  http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spy := &spyServiceRunner{}
			srv := &server{
				router:  tt.fields.router,
				logger:  tt.fields.logger,
				service: spy,
			}

			srv.Routes()

			req := httptest.NewRequest(http.MethodPost, "/v1/services/nginx", tt.body)
			w := httptest.NewRecorder()

			srv.ServeHTTP(w, req)

			if w.Result().StatusCode != tt.statusCode {
				t.Errorf("expected status code %v, got %v instead", tt.statusCode, w.Result().StatusCode)
			}
			if tt.wantCommand != "" {
				if spy.Command != tt.wantCommand {
					t.Errorf("wanted the command %v, got %v instead", tt.wantCommand, spy.Command)
				}
			}
			if len(tt.wantArgs) > 0 {
				if !reflect.DeepEqual(spy.Args, tt.wantArgs) {
					t.Errorf("expected the args to be the same")
				}
			}
		})
	}
}
