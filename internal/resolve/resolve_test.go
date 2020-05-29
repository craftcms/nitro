package resolve

import (
	"os"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func TestAbsPath(t *testing.T) {
	home, err := homedir.Dir()
	if err != nil {
		t.Fatal(err)
	}
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		path string
		home string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "can get the full path",
			args:    args{path: "./testdata/exists", home: home},
			want:    dir + "/testdata/exists",
			wantErr: false,
		},
		{
			name:    "can resolve tilde path references",
			args:    args{path: "~/go/src/github.com/craftcms/nitro/internal/resolve", home: home},
			want:    dir,
			wantErr: false,
		},
		{
			name:    "non-existent directories return errors",
			args:    args{path: "./testdata/empty", home: home},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AbsPath(tt.args.path, tt.args.home)
			if (err != nil) != tt.wantErr {
				t.Errorf("AbsPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AbsPath() \ngot = \n%v \nwant = \n%v", got, tt.want)
			}
		})
	}
}
