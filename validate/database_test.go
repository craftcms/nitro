package validate

import (
	"testing"

	"github.com/craftcms/nitro/config"
)

func TestDatabaseEngine(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "mysql does not return error",
			args:    args{v: "mysql"},
			wantErr: false,
		},
		{
			name:    "mySQL returns error",
			args:    args{v: "mySQL"},
			wantErr: true,
		},
		{
			name:    "postgres does not return error",
			args:    args{v: "postgres"},
			wantErr: false,
		},
		{
			name:    "postgresql returns error",
			args:    args{v: "postgresql"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DatabaseEngine(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("DatabaseEngine() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


func TestDatabaseVersion(t *testing.T) {
	type args struct {
		e string
		v string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "postgres and supported version 12.2 does not return error",
			args: args{
				e: "postgres",
				v: "12.2",
			},
			wantErr: false,
		},
		{
			name: "postgres and supported version 12 does not return error",
			args: args{
				e: "postgres",
				v: "12",
			},
			wantErr: false,
		},
		{
			name: "postgres and supported version 11.7 does not return error",
			args: args{
				e: "postgres",
				v: "11.7",
			},
			wantErr: false,
		},
		{
			name: "postgres and supported version 11 does not return error",
			args: args{
				e: "postgres",
				v: "11",
			},
			wantErr: false,
		},
		{
			name: "postgres and supported version 9.5 does not return error",
			args: args{
				e: "postgres",
				v: "9.5",
			},
			wantErr: false,
		},
		{
			name: "postgres and supported version 9 does not return error",
			args: args{
				e: "postgres",
				v: "9",
			},
			wantErr: false,
		},
		{
			name: "postgres and supported version 9.6 does not return error",
			args: args{
				e: "postgres",
				v: "9.6",
			},
			wantErr: false,
		},
		{
			name: "postgres and supported version 10 does not return error",
			args: args{
				e: "postgres",
				v: "10",
			},
			wantErr: false,
		},
		{
			name: "unsupported engine returns error",
			args: args{
				e: "notsupported",
				v: "1.0",
			},
			wantErr: true,
		},
		{
			name: "supported engine and version does not return error",
			args: args{
				e: "mysql",
				v: "5.7",
			},
			wantErr: false,
		},
		{
			name: "supported engine and version does not return error",
			args: args{
				e: "mysql",
				v: "8.0",
			},
			wantErr: false,
		},
		{
			name: "supported engine and version does not return error",
			args: args{
				e: "mysql",
				v: "8",
			},
			wantErr: false,
		},
		{
			name: "supported engine and version does not return error",
			args: args{
				e: "mysql",
				v: "5.8",
			},
			wantErr: false,
		},
		{
			name: "supported engine and version does not return error",
			args: args{
				e: "mysql",
				v: "5",
			},
			wantErr: false,
		},
		{
			name: "supported engine and unsupported version returns error",
			args: args{
				e: "mysql",
				v: "5.1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DatabaseEngineAndVersion(tt.args.e, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("DatabaseEngineAndVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseConfig(t *testing.T) {
	type args struct {
		databases []config.Database
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "non duplicate ports does not return an error",
			args: args{
				databases: []config.Database{
					{
						Engine:  "mysql",
						Version: "8",
						Port:    "3306",
					},
					{
						Engine:  "mysql",
						Version: "5.7",
						Port:    "33061",
					},
					{
						Engine:  "postgres",
						Version: "11",
						Port:    "5432",
					},
					{
						Engine:  "postgres",
						Version: "12",
						Port:    "54321",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "duplicate ports return error",
			args: args{
				databases: []config.Database{
					{
						Engine:  "mysql",
						Version: "5.7",
						Port:    "3306",
					},
					{
						Engine:  "mysql",
						Version: "8",
						Port:    "3306",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "duplicate engines and version returns an error",
			args: args{
				databases: []config.Database{
					{
						Engine:  "mysql",
						Version: "5.7",
						Port:    "3306",
					},
					{
						Engine:  "mysql",
						Version: "5.7",
						Port:    "33061",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DatabaseConfig(tt.args.databases); (err != nil) != tt.wantErr {
				t.Errorf("DatabaseConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseName(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid database names do not return an error",
			args: args{
				s: "this_is_a_valid_9_name",
			},
			wantErr: false,
		},
		{
			name: "longer than 64 chars returns an error",
			args: args{
				s: "A5C35D63132B33B572C6E8A7A24B4BE875015AC21AF972D0FF8F3D1396680D878",
			},
			wantErr: true,
		},
		{
			name: "empty strings returns an error",
			args: args{
				s: "",
			},
			wantErr: true,
		},
		{
			name: "dollar signs return an error",
			args: args{
				s: "$1234",
			},
			wantErr: true,
		},
		{
			name: "dashes return an error",
			args: args{
				s: "test-database-name",
			},
			wantErr: true,
		},
		{
			name: "cannot start with pg_",
			args: args{
				s: "pg_databasename",
			},
			wantErr: true,
		},
		{
			name: "numbers in the front return error",
			args: args{
				s: "9this",
			},
			wantErr: true,
		},
		{
			name: "spaces return error",
			args: args{
				s: " this is a space",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DatabaseName(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("DatabaseName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
