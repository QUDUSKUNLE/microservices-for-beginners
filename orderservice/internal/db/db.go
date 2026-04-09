package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "order.db")
	if err != nil {
		return nil, err
	}
	if err = Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS orders(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_email TEXT,
		product_id INTEGER,
		quantity INTEGER,
		address TEXT,
		status TEXT
	);`
	_, err := db.Exec(schema)
	return err
}
