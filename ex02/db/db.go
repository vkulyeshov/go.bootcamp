package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	createTable := `
	CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		link TEXT
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func SaveItem(db *sql.DB, title, link string) error {
	_, err := db.Exec("INSERT INTO items (title, link) VALUES (?, ?)", title, link)
	return err
}

func ClearItems(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM items;")
	return err
}
