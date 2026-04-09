package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	// basic connection check
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// create table
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		email TEXT PRIMARY KEY,
		passwordhash TEXT,
		address TEXT
	);
	`

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	log.Println("DB initialized")

	return db, nil
}
