package database

import (
	"os"
	"testing"
)

func TestDetermineEngine(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "can detect mysql database backup files",
			args:    args{file: "./testdata/mysql-backup.sql"},
			want:    "mysql",
			wantErr: false,
		},
		{
			name:    "can detect postgres database backup files",
			args:    args{file: "./testdata/postgres-backup.sql"},
			want:    "postgres",
			wantErr: false,
		},
		{
			name:    "non mysql or postgres files return an error",
			args:    args{file: "./testdata/random.txt"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.args.file)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			got, err := DetermineEngine(f)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetermineEngine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DetermineEngine() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasCreateStatement(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "can detect when mysql files have a create database statement",
			args:    args{file: "./testdata/mysql-create-backup.sql"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "can detect when mysql files have a create database statement example two",
			args:    args{file: "./testdata/mysql-backup-example-two.sql"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "can detect when mysql files does not have a create database statement",
			args:    args{file: "./testdata/mysql-backup.sql"},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.args.file)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			got, err := HasCreateStatement(f)
			if (err != nil) != tt.wantErr {
				t.Errorf("HasCreateStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HasCreateStatement() got = %v, want %v", got, tt.want)
			}
		})
	}
}
