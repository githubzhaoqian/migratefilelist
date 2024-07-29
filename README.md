# migratefilelist

golang-migrate/migrate 的 source 扩展，解决多版本同时开发并不确定合并顺序的问题，用文件行号代替原本文件的版本号。

golang-migrate/migrate source extension，Solves the problem of simultaneous development of multiple versions without determining the merge order，Replace the version number of the original file with the file line number.

# Example

migrate.list
```
20240729_show.up.sql
20240729_show2.up.sql
```
> 不能有空行和相同的行
>
> can't have empty rows and the same rows

run.go
```go

package main

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

func main() {
	db, _ := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/migrate_test?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai")
	driver, _ := mysql.WithInstance(db, &mysql.Config{MigrationsTable: "schema_migrations_order"})
	m, _ := migrate.NewWithDatabaseInstance(
		"filelist://migrate.list",
		"mysql",
		driver,
	)
	version, dirty, err := m.Version()
	if err != nil {
		panic(err)
	}
	fmt.Println("current versio: %d, dirty: %t", version, dirty)
	m.Steps(2)
}
```

