package keys

import (
	"reflect"
	"testing"
)

func TestFind(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "can find all of the keys",
			args:    args{path: "./testdata"},
			want:    map[string]string{"id_rsa": "id_rsa.pub"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Find(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = \n%v, \nwant \n%v", got, tt.want)
			}
		})
	}
}
