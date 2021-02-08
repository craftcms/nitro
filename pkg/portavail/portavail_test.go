package portavail

import (
	"fmt"
	"net"
	"strconv"
	"testing"
)

func TestCheck(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Logf("error creating listener, %s", err.Error())
		t.Fail()
	}
	p := lis.Addr().(*net.TCPAddr).Port
	usedPort := strconv.Itoa(p)

	type args struct {
		host string
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
			if err := Check(tt.args.host, tt.args.port); (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindNext(t *testing.T) {
	// find a random port
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Logf("error creating listener, %s", err.Error())
		t.Fail()
	}
	p := lis.Addr().(*net.TCPAddr).Port
	usedPort := strconv.Itoa(p)

	type args struct {
		host string
		port string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "find the next available port",
			args: args{
				port: usedPort,
			},
			want:    fmt.Sprintf("%d", p+1),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindNext(tt.args.host, tt.args.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindNext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FindNext() = %v, want %v", got, tt.want)
			}
		})
	}
}
