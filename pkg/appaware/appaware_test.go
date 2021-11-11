package appaware

import (
	"reflect"
	"testing"

	config "github.com/craftcms/nitro/pkg/config/v3"
)

func TestDetect(t *testing.T) {
	type args struct {
		cfg config.Config
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    *config.App
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Detect(tt.args.cfg, tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("Detect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Detect() got = %v, want %v", got, tt.want)
			}
		})
	}
}
