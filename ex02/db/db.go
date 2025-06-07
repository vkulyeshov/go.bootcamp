package db

import (
	"database/sql"
	"rss_fetcher/data"

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
		link TEXT,
		description TEXT,
		author TEXT,
		category TEXT,
		pub_date TEXT,	
		guid TEXT	
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func SaveItem(db *sql.DB, title, link, description, author, category, pub_date, guid string) error {
	_, err := db.Exec("INSERT INTO items (title, link, description, author, category, pub_date, guid) VALUES (?, ?, ?, ?, ?, ?, ?)",
		title, link, description, author, category, pub_date, guid)
	return err
}

func ClearItems(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM items;")
	return err
}

func LoadItems(db *sql.DB) ([]data.RSSItem, error) {
	rows, err := db.Query("SELECT id, title, link, description, author, category, pub_date, guid FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]data.RSSItem, 0)

	for rows.Next() {
		var id int
		var title, link, description, author, category, pub_date, guid string

		if err := rows.Scan(&id, &title, &link, &description, &author, &category, &pub_date, &guid); err != nil {
			return nil, err
		}

		item := data.RSSItem{
			ID:          id,
			Title:       title,
			Link:        link,
			Description: description,
			Author:      author,
			Category:    category,
			PubDate:     pub_date,
			GUID:        guid,
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
