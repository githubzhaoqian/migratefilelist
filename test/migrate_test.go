package test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"

	_ "github.com/go-sql-driver/mysql"

	"github.com/githubzhaoqian/migrate-filelist/source/iofs"

	"github.com/golang-migrate/migrate/v4/database/mysql"

	_ "github.com/githubzhaoqian/migrate-filelist/source/filelist"
)

func TestRun(t *testing.T) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/migrate_test?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai")
	if err != nil {
		t.Fatal(err)
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{MigrationsTable: "schema_migrations_order"})
	if err != nil {
		t.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"filelist://migrate.list",
		"mysql",
		driver,
	)
	version, dirty, err := m.Version()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("current versio: %d, dirty: %t", version, dirty)
	m.Steps(2)
	t.Log("END")
}

func TestParse(t *testing.T) {
	m, err := source.Parse("123_name.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Migration: ", m)
	t.Log("END")
}

func TestNew(t *testing.T) {
	driver, err := iofs.New("migrate.list", ".")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("driver: ", driver)
	t.Log("END")
}

func TestFileList(t *testing.T) {
	schemaFile, err := os.Open("migrate.list")
	if err != nil {
		t.Fatal(err)
	}
	fInfo, err := schemaFile.Stat()
	if err != nil {
		t.Fatal(err)
	}
	err = schemaFile.Close()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("fileName", fInfo.Name())
	t.Log("END")
}
