package config

import (
	"testing"
)

func TestDatabase_GetHostname(t *testing.T) {
	type fields struct {
		Engine    string
		Version   string
		Port      string
		Ephemeral bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name:    "can get the hostname for a database container",
			fields:  fields{Engine: "mysql", Version: "5.7", Port: "3306"},
			want:    "mysql-5.7-3306",
			wantErr: false,
		},
		{
			name:    "empty values return an error",
			fields:  fields{Engine: "mysql", Port: "3306"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				Engine:  tt.fields.Engine,
				Version: tt.fields.Version,
				Port:    tt.fields.Port,
			}
			got, err := d.GetHostname()
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.GetHostname() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Database.GetHostname() = %v, want %v", got, tt.want)
			}
		})
	}
}
