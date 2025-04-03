package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/pressly/goose/v3"
)

func main() {
	var migrationsPath, storagePath string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations folder")
	flag.StringVar(&storagePath, "storage-path", "", "path to storage file")
	flag.Parse()

	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	if storagePath == "" {
		panic("storage-path is required")
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", storagePath))
	if err != nil {
		panic(err)
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		panic(err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
}
