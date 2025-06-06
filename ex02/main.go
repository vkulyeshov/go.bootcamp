package main

import (
	"flag"
	"fmt"
	"log"
	"rss_fetcher/db"
	"rss_fetcher/parser"
)

func main() {
	rssURL := flag.String("url", "", "RSS feed URL (required)")
	dbPath := flag.String("db", "rss_items.db", "Path to SQLite database")
	limit := flag.Int("limit", 0, "Max number of entries (0 = all)")
	reset := flag.Bool("reset", false, "Clear table before inserting")
	flag.Parse()

	if *rssURL == "" {
		log.Fatal("The --url parameter is required")
	}

	feedItems, err := parser.FetchRSS(*rssURL)
	if err != nil {
		log.Fatalf("Error fetching RSS: %v", err)
	}

	dbConn, err := db.InitDB(*dbPath)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer dbConn.Close()

	if *reset {
		if err := db.ClearItems(dbConn); err != nil {
			log.Fatalf("Error clearing table: %v", err)
		}
	}

	count := 0
	for _, item := range feedItems {
		if *limit > 0 && count >= *limit {
			break
		}
		err := db.SaveItem(dbConn, item.Title, item.Link)
		if err != nil {
			log.Printf("Insert error: %v", err)
		} else {
			count++
			log.Printf("record #: %d", count)
		}
	}

	fmt.Printf("Saved %d items.", count)
}
