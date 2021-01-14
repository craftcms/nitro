package setup

import (
	"bytes"
	"io"
	"testing"
)

func Test_setupPostgres(t *testing.T) {
	type args struct {
		rdr      io.Reader
		msg      string
		fallback bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "returns true if yes is entered",
			args: args{
				rdr: bytes.NewReader([]byte("Y\n")),
				msg: "some message",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false if no is entered",
			args: args{
				rdr: bytes.NewReader([]byte("\n")),
				msg: "some message",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := confirm(tt.args.rdr, tt.args.msg, tt.args.fallback)
			if (err != nil) != tt.wantErr {
				t.Errorf("setupPostgres() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("setupPostgres() = %v, want %v", got, tt.want)
			}
		})
	}
}
