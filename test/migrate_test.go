package test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"

	_ "github.com/go-sql-driver/mysql"

	"github.com/githubzhaoqian/migratefilelist/source/iofs"

	"github.com/golang-migrate/migrate/v4/database/mysql"

	_ "github.com/githubzhaoqian/migratefilelist/source/filelist"
)

func TestRun(t *testing.T) {
	db, err := sql.Open("mysql", "root:gogogo@tcp(127.0.0.1:3306)/gofast?collation=utf8mb4_bin&parseTime=true&loc=Local&multiStatements=true")
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
	m.Steps(2)
	version, dirty, err := m.Version()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("current versio: %d, dirty: %t", version, dirty)
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
