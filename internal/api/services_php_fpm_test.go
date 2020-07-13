package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_server_handlePhpFpmService(t *testing.T) {
	type fields struct {
		router *http.ServeMux
		logger *log.Logger
	}
	tests := []struct {
		name       string
		fields     fields
		statusCode int
	}{
		{
			name: "testing request validation",
			fields: fields{
				router: http.NewServeMux(),
				logger: log.New(ioutil.Discard, "nitrod test ", 0),
			},
			statusCode: http.StatusUnprocessableEntity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &server{
				router:  tt.fields.router,
				logger:  tt.fields.logger,
				service: spyServiceRunner{},
			}

			srv.Routes()

			req := httptest.NewRequest(http.MethodPost, "/v1/services/php-fpm", strings.NewReader(""))
			w := httptest.NewRecorder()

			srv.ServeHTTP(w, req)

			if w.Result().StatusCode != tt.statusCode {
				t.Errorf("expected status code %v, got %v instead", tt.statusCode, w.Result().StatusCode)
			}
		})
	}
}

type spyServiceRunner struct {
	Command string
	Args    []string
}

func (r spyServiceRunner) Run(command string, args []string) ([]byte, error) {
	r.Command = command
	r.Args = args

	return []byte("test"), nil
}
