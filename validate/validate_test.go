package validate

import (
	"testing"

	"github.com/craftcms/nitro/config"
)

func TestPHPVersionFlag(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "7.4",
			args: args{
				"7.4",
			},
			wantErr: false,
		},
		{
			name: "7.3",
			args: args{
				"7.3",
			},
			wantErr: false,
		},
		{
			name: "7.2",
			args: args{
				"7.2",
			},
			wantErr: false,
		},
		{
			name: "junk",
			args: args{
				"junk",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PHPVersion(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("PHPVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseFlag(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "mariadb",
			args:    args{v: "mariadb"},
			wantErr: false,
		},
		{
			name:    "maria",
			args:    args{v: "maria"},
			wantErr: true,
		},
		{
			name:    "postgres",
			args:    args{v: "postgres"},
			wantErr: false,
		},
		{
			name:    "postgresql",
			args:    args{v: "postgresql"},
			wantErr: true,
		},
		{
			name:    "pgsql",
			args:    args{v: "pgsql"},
			wantErr: true,
		},
		{
			name:    "cockroach",
			args:    args{v: "cockroach"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Database(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Database() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDomain(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "is valid",
			args:    args{v: "example.test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Domain(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Domain() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "integers only return error",
			args:    args{v: "2"},
			wantErr: true,
		},
		{
			name:    "lower case G return error",
			args:    args{v: "2g"},
			wantErr: true,
		},
		{
			name:    "proper values do not return error",
			args:    args{v: "2G"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Memory(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Memory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
