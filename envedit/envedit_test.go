package envedit

import "testing"

func TestEdit(t *testing.T) {
	type args struct {
		file    string
		updates []map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "files are updated",
			args: args{
				file: "testdata/env-example",
				updates: []map[string]string{
					{
						"DB_SERVER": "dddd",
					},
				},
			},
		},
		{
			name:    "no file returns an error",
			args:    args{file: "testdata/nothing-here"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Edit(tt.args.file, tt.args.updates...); (err != nil) != tt.wantErr {
				t.Errorf("Edit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
