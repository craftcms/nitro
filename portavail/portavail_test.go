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
		port string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "used ports return an error",
			args:    args{port: usedPort},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Check(tt.args.port); (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
