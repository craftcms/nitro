package envedit

import (
	"testing"
)

func TestEdit(t *testing.T) {
	type args struct {
		file    string
		updates map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "environment variables are updated",
			args: args{
				file: "testdata/env-example",
				updates: map[string]string{
					"DB_DRIVER":   "pgsql",
					"DB_SERVER":   "postgres-13-5432",
					"DB_PORT":     "5432",
					"DB_DATABASE": "example",
					"DB_USER":     "nitro",
					"DB_PASSWORD": "nitro",
				},
			},
			want: `# The environment Craft is currently running in (dev, staging, production, etc.)
ENVIRONMENT=dev

# The application ID used to to uniquely store session and cache data, mutex locks, and more
APP_ID=CraftCMS

# The secure key Craft will use for hashing and encrypting data
SECURITY_KEY=

# The database driver that will be used (mysql or pgsql)
DB_DRIVER=pgsql

# The database server name or IP address
DB_SERVER=postgres-13-5432

# The port to connect to the database with
DB_PORT=5432

# The name of the database to select
DB_DATABASE=example

# The database username to connect with
DB_USER=nitro

# The database password to connect with
DB_PASSWORD=nitro

# The database schema that will be used (PostgreSQL only)
DB_SCHEMA=public

# The prefix that should be added to generated table names (only necessary if multiple things are sharing the same database)
DB_TABLE_PREFIX=
`,
			wantErr: false,
		},
		{
			name:    "missing file returns error",
			args:    args{file: "testdata/empty"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Edit(tt.args.file, tt.args.updates)
			if (err != nil) != tt.wantErr {
				t.Errorf("Edit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Edit() = %v, want %v", got, tt.want)
			}
		})
	}
}
