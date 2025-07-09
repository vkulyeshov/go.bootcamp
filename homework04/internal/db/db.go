package db

import (
	"context"
	"errors"
	"rss_fetcher/internal/data"
	"time"

	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type QueryInterface interface {
	Query(sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(sql string, args ...interface{}) pgx.Row
	Exec(sql string, args ...interface{}) (pgconn.CommandTag, error)
}

type ConnectionQuery struct {
	db *pgx.Conn
}

type TransactionQuery struct {
	tx pgx.Tx
}

func NewTransactionQuery(tx pgx.Tx) *TransactionQuery {
	return &TransactionQuery{tx: tx}
}

func NewConnectionQuery(db *pgx.Conn) *ConnectionQuery {
	return &ConnectionQuery{db: db}
}

func (c *ConnectionQuery) Query(sql string, args ...interface{}) (pgx.Rows, error) {
	return c.db.Query(context.Background(), sql, args...)
}

func (c *ConnectionQuery) QueryRow(sql string, args ...interface{}) pgx.Row {
	return c.db.QueryRow(context.Background(), sql, args...)
}

func (c *ConnectionQuery) Exec(sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return c.db.Exec(context.Background(), sql, args...)
}

func (t *TransactionQuery) Exec(sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return t.tx.Exec(context.Background(), sql, args...)
}

func (t *TransactionQuery) Query(sql string, args ...interface{}) (pgx.Rows, error) {
	return t.tx.Query(context.Background(), sql, args...)
}

func (t *TransactionQuery) QueryRow(sql string, args ...interface{}) pgx.Row {
	return t.tx.QueryRow(context.Background(), sql, args...)
}

func InitDB(connectionUrl string) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), connectionUrl)
}

func StartTransaction(db *pgx.Conn) (pgx.Tx, error) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func CommitTransaction(tx pgx.Tx) error {
	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}

func RollbackTransaction(tx pgx.Tx) error {
	if err := tx.Rollback(context.Background()); err != nil {
		return err
	}
	return nil
}

func AddJob(db QueryInterface, link string) error {
	_, err := db.Exec("INSERT INTO news_jobs (link) VALUES ($1)", link)
	return err
}

func DeleteJob(db QueryInterface, id int) error {
	_, err := db.Exec("DELETE FROM news_jobs WHERE job_id = $1", id)
	return err
}

func LoadJobs(db QueryInterface, filter string) ([]data.NewsJob, error) {
	if filter != "" {
		filter = " WHERE " + filter
	}
	rows, err := db.Query("SELECT job_id, link, status, created_at, updated_at FROM news_jobs" + filter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := make([]data.NewsJob, 0)

	for rows.Next() {
		var id, status int
		var link string
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&id, &link, &status, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		job := data.NewsJob{
			ID:        id,
			Link:      link,
			Status:    status,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func UpdateJobStatus(db QueryInterface, id int, status int) error {
	_, err := db.Exec("UPDATE news_jobs SET status = $1, updated_at = NOW() WHERE job_id = $2", status, id)
	return err
}

func AddChannel(db *pgx.Conn, link string) error {
	_, err := db.Exec(context.Background(), "INSERT INTO channels (link) VALUES ($1)", link)
	return err
}

func DeleteChannel(db *pgx.Conn, id int) error {
	_, err := db.Exec(context.Background(), "DELETE FROM channels WHERE channel_id = $1", id)
	return err
}

func DeleteChannels(db *pgx.Conn) error {
	_, err := db.Exec(context.Background(), "DELETE FROM channels;")
	return err
}

func LoadChannels(db *pgx.Conn) ([]data.Channel, error) {
	rows, err := db.Query(context.Background(), "SELECT channel_id, link, title, description, rss_link, last_updated FROM channels")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	channels := make([]data.Channel, 0)

	for rows.Next() {
		var id int
		var link, title, description, rssLink string
		var updatedAt time.Time

		if err := rows.Scan(&id, &link, &title, &description, &rssLink, &updatedAt); err != nil {
			return nil, err
		}

		channel := data.Channel{
			ID:          id,
			Link:        link,
			Title:       title,
			Description: description,
			RSSLink:     rssLink,
			UpdatedAt:   updatedAt,
		}
		channels = append(channels, channel)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return channels, nil
}

func LoadChannel(db *pgx.Conn, id int) (*data.Channel, error) {
	row, err := db.Query(context.Background(), "SELECT link, title, rss_link, description, last_updated FROM channels WHERE channel_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	if !row.Next() {
		return nil, errors.New("channel not found")
	}
	var title, link, rssLink, description string
	var lastUpdated time.Time
	if err := row.Scan(&id, &title, &link, &rssLink, &description, &lastUpdated); err != nil {
		return nil, err
	}

	result := &data.Channel{
		ID:          id,
		Title:       title,
		Link:        link,
		RSSLink:     rssLink,
		Description: description,
		UpdatedAt:   lastUpdated,
	}

	items, err := loadNews(&ConnectionQuery{db}, "channel_id = "+fmt.Sprint(result.ID))
	if err != nil {
		return nil, err
	}
	result.Items = items
	return result, nil
}

func AddNews(db QueryInterface, channelID int, news data.ChannelNews) error {
	_, err := db.Exec("INSERT INTO channel_news (channel_id, title, link, description, author, category, pub_date, guid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		channelID, news.Title, news.Link, news.Description, news.Author, news.Category, news.PubDate, news.GUID)
	return err
}

func DeleteNews(db *pgx.Conn, id int) error {
	_, err := db.Exec(context.Background(), "DELETE FROM channel_news WHERE news_id = $1", id)
	return err
}

func loadNews(db QueryInterface, filter string) ([]data.ChannelNews, error) {
	result := make([]data.ChannelNews, 0)

	var query string = "SELECT news_id, title, link, description, author, category, pub_date, guid FROM channel_news"
	if filter != "" {
		query += " WHERE " + filter
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var title, link, description, author, category, guid string
		var pubDate time.Time

		if err := rows.Scan(&id, &title, &link, &description, &author, &category, &pubDate, &guid); err != nil {
			return nil, err
		}

		item := data.ChannelNews{
			ID:          id,
			Title:       title,
			Link:        link,
			Description: description,
			Author:      author,
			Category:    category,
			PubDate:     pubDate,
			GUID:        guid,
		}
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func LoadNews(db QueryInterface) ([]data.ChannelNews, error) {
	return loadNews(db, "")
}

func LoadNewsByLink(db QueryInterface, link string) ([]data.ChannelNews, error) {
	return loadNews(db, "link = '"+link+"'")
}

func Close(db *pgx.Conn) {
	if db != nil {
		db.Close(context.Background())
	}
}
