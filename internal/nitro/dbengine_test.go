package nitro

import (
	"path/filepath"
	"testing"
)

func TestFindEngineByDump(t *testing.T) {
	t.Run("Mysql engine is detected", func(t *testing.T) {
		mySqlEngine := "mysql"
		mySqlFileAbsPath, err := filepath.Abs("./testdata/dbengine/mysqldump.sql")
		if err != nil {
			t.Errorf("Cannot find MySQL test file for DB engine detection")
			return
		}

		engine := FindEngineByDump(mySqlFileAbsPath)
		if engine != mySqlEngine {
			t.Errorf("engine returned as \"%v\", want \"%v\"", engine, mySqlEngine)
			return
		}
	})

	t.Run("Postgres engine is detected", func(t *testing.T) {
		postgresEngine := "postgres"
		mySqlFileAbsPath, err := filepath.Abs("./testdata/dbengine/postgresdump.sql")
		if err != nil {
			t.Errorf("Cannot find Postgres test file for DB engine detection")
			return
		}

		engine := FindEngineByDump(mySqlFileAbsPath)
		if engine != postgresEngine {
			t.Errorf("engine returned as \"%v\", want \"%v\"", engine, postgresEngine)
			return
		}
	})
}
