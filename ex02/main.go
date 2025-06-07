package main

import (
	"flag"
	"fmt"
	"log"
	"rss_fetcher/db"
	"rss_fetcher/parser"
)

const SPLIT_LINE = "==========================================="

func main() {
	rssURL := flag.String("url", "", "RSS feed URL")
	dbPath := flag.String("db", "rss_items.db", "Path to SQLite database")
	limit := flag.Int("limit", 0, "Max number of entries (0 = all)")
	reset := flag.Bool("reset", false, "Clear table before inserting")
	showDB := flag.Bool("show-db", false, "Show table contents")
	flag.Parse()

	if *rssURL == "" && !*showDB && !*reset {
		log.Fatal("Please specify either --url or --show-db or reset parameter")
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
		log.Println("Database is cleaned!")
	}

	if *rssURL != "" {
		items, err := parser.FetchRSS(*rssURL)
		if err != nil {
			log.Fatalf("Error fetching RSS: %v", err)
		}

		count := 0
		for _, item := range items {
			if *limit > 0 && count >= *limit {
				break
			}
			log.Printf("Storing: %s ...", item.Title)

			err := db.SaveItem(dbConn, item.Title, item.Link, item.Description, item.Author, item.Category, item.PubDate, item.GUID)

			if err != nil {
				log.Printf("Insert error: %v", err)
			} else {
				count++
			}
		}
		log.Printf("Stored %d items to database", count)
	}

	if *showDB {
		items, err := db.LoadItems(dbConn)
		if err != nil {
			log.Fatalf("Error loading items from database: %v", err)
		}

		fmt.Println(SPLIT_LINE)
		fmt.Println("ID | Title | Link")
		fmt.Println(SPLIT_LINE)

		count := 0
		for _, item := range items {
			fmt.Printf("%d | %s | %s\n", item.ID, item.Title, item.Link)
			count++
		}

		log.Println(SPLIT_LINE)
		log.Println("== Found ", count, "items in the database")
		log.Println(SPLIT_LINE)
	}
}
