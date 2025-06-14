package db

import (
	"context"

	"rss_fetcher/data"

	"github.com/jackc/pgx/v5"
)

func InitDB(connectionUrl string) (*pgx.Conn, error) {
	db, err := pgx.Connect(context.Background(), connectionUrl)
	if err != nil {
		return nil, err
	}
	createTable := `
	CREATE TABLE IF NOT EXISTS items (
		id SERIAL PRIMARY KEY,
		title TEXT,
		link TEXT,
		description TEXT,
		author TEXT,
		category TEXT,
		pub_date TEXT,	
		guid TEXT	
	);`
	_, err = db.Exec(context.Background(), createTable)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func SaveItem(db *pgx.Conn, title, link, description, author, category, pub_date, guid string) error {
	_, err := db.Exec(context.Background(), "INSERT INTO items (title, link, description, author, category, pub_date, guid) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		title, link, description, author, category, pub_date, guid)
	return err
}

func ClearItems(db *pgx.Conn) error {
	_, err := db.Exec(context.Background(), "DELETE FROM items;")
	return err
}

func LoadItems(db *pgx.Conn) ([]data.RSSItem, error) {
	rows, err := db.Query(context.Background(), "SELECT id, title, link, description, author, category, pub_date, guid FROM items")
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

func Close(db *pgx.Conn) {
	if db != nil {
		db.Close(context.Background())
	}
}
