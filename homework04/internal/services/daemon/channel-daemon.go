package daemon

import (
	"context"
	"fmt"
	"log"
	"rss_fetcher/internal/data"
	"rss_fetcher/internal/db"
	"rss_fetcher/internal/parser"

	"github.com/jackc/pgx/v5"
)

type ChannelDaemon struct {
	dbConn *pgx.Conn
}

func NewChannelDaemon(dbConn *pgx.Conn) *ChannelDaemon {
	return &ChannelDaemon{dbConn: dbConn}
}

func (daemon *ChannelDaemon) CheckFeeds() {
	channels, err := db.LoadChannels(daemon.dbConn)
	if err != nil {
		log.Printf("Failed to load channels: %v", err)
		return
	}

	for _, channel := range channels {
		if err := daemon.processChannel(channel); err != nil {
			log.Printf("Error processing channel %s: %v", channel.Link, err)
		}
	}
}

func (daemon *ChannelDaemon) saveNews(channelID int, news data.ChannelNews) error {
	tx, err := db.StartTransaction(daemon.dbConn)
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return err
	}
	defer tx.Rollback(context.Background())

	transactionQuery := db.NewTransactionQuery(tx)

	if rows, err := db.LoadNewsByLink(transactionQuery, news.Link); err != nil {
		log.Printf("Failed to load news by link %s: %v", news.Link, err)
		return err
	} else if len(rows) > 0 {
		log.Printf("News item %s already exists, skipping", news.Link)
		return nil
	}
	log.Printf("Found news: %s", news.Link)
	if err := db.AddNews(transactionQuery, channelID, news); err != nil {
		log.Printf("Failed to save news item %s: %v", news.Link, err)
		return err
	}

	if db.AddJob(transactionQuery, news.Link) != nil {
		log.Printf("Failed to add job for news item %s: %v", news.Link, err)
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("commit error: %w", err)
	}

	return nil
}

func (daemon *ChannelDaemon) processChannel(channel data.Channel) error {
	news, err := parser.FetchRSS(channel.Link)
	if err != nil {
		log.Printf("Failed to fetch RSS for channel %s: %v", channel.Link, err)
		return err
	}
	for _, item := range news {
		daemon.saveNews(channel.ID, item)
		if err != nil {
			log.Printf("Failed to save news item %s: %v", item.Link, err)
		}
	}
	log.Printf("Processed channel %s with %d items", channel.Link, len(news))
	return nil
}
