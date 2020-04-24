package config

import "testing"

func TestMount_Exists(t *testing.T) {
	type fields struct {
		Source string
		Dest   string
	}
	type args struct {
		webroot string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "returns true if the sites mount exists",
			fields: fields{
				Source: "./testdata/example-source",
				Dest:   "/nitro/sites/example",
			},
			args: args{
				webroot: "/nitro/sites/example/web",
			},
			want: true,
		},
		{
			name: "returns false if the sites mount does not exist",
			fields: fields{
				Source: "./testdata/another-source",
				Dest:   "/nitro/sites/another",
			},
			args: args{
				webroot: "/nitro/sites/example/web",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mount{
				Source: tt.fields.Source,
				Dest:   tt.fields.Dest,
			}
			if got := m.Exists(tt.args.webroot); got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}
