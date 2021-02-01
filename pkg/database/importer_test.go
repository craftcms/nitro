package database

import (
	"os/exec"
	"testing"
)

func TestDefaultImportToolFinder(t *testing.T) {
	shPath, err := exec.LookPath("sh")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		engine           string
		version          string
		postgresToolPath string
		mysqlToolPath    string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "mysql engine and existing path returns an error",
			args: args{
				engine:        "mysql",
				version:       "",
				mysqlToolPath: shPath,
			},
			want:    shPath,
			wantErr: false,
		},
		{
			name: "postgres engine and existing path returns an error",
			args: args{
				engine:           "postgres",
				version:          "",
				postgresToolPath: shPath,
			},
			want:    shPath,
			wantErr: false,
		},
		{
			name: "postgres engine and missing path returns an error",
			args: args{
				engine:           "postgres",
				version:          "",
				postgresToolPath: "missingpath",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "mysql engine and missing path returns an error",
			args: args{
				engine:        "mysql",
				version:       "",
				mysqlToolPath: "missingpath",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "no matching engine or version returns an error",
			args: args{
				engine:  "",
				version: "",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if tt.args.mysqlToolPath != "" {
			MySQLImportCommand = tt.args.mysqlToolPath
		}
		if tt.args.postgresToolPath != "" {
			PostgresImportCommand = tt.args.postgresToolPath
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := DefaultImportToolFinder(tt.args.engine, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultImportToolFinder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DefaultImportToolFinder() = %v, want %v", got, tt.want)
			}
		})
	}
}
