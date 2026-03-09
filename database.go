package main

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"os"
	"path"

	"sensor_dashboard/db"

	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var ddl string

func InitDB(datasourceDir string) *db.Queries {
	ctx := context.Background()

	_, err := os.ReadDir(datasourceDir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(datasourceDir, 0755)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}

	sqliteDB, err := sql.Open("sqlite3", path.Join(datasourceDir, "sqlite.db"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = sqliteDB.ExecContext(ctx, ddl)
	if err != nil {
		log.Fatal(err)
	}

	queries := db.New(sqliteDB)

	return queries
}
