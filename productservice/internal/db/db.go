package db

import (
	 _ "modernc.org/sqlite"
	"database/sql"

)


func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "products.db")
	if err != nil {
		return nil, err
	}

	err = Migrate(db)
	if err != nil {
		return nil, err
	}
	return db, nil

}

func Migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS  products(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	price REAL,
	category TEXT,
	stock INTEGER
	);`
	_, err := db.Exec(schema)
	return err
}
