package portavail

import (
	"net"
	"strconv"
	"testing"
)

func TestCheck(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	p := lis.Addr().(*net.TCPAddr).Port
	usedPort := strconv.Itoa(p)

	type args struct {
		ports []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "returns error when no ports are provided",
			wantErr: true,
		},
		{
			name:    "used ports return an error",
			args:    args{ports: []string{usedPort}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Check(tt.args.ports...); (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
