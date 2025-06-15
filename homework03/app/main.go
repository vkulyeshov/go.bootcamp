package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"rss_fetcher/db"
	"rss_fetcher/parser"

	"github.com/joho/godotenv"
)

const SPLIT_LINE = "==========================================="

func main() {
	// Artificial solution should be replaced with a proper config management
	defaultConnection := "postgres://rss:rss@localhost:5432/rss"
	err := godotenv.Load("../.env")
	if err == nil {
		defaultConnection =
			fmt.Sprintf("postgres://%s:%v@%s:%s/%s",
				os.Getenv("POSTGRES_USER"),
				os.Getenv("POSTGRES_PASSWORD"),
				os.Getenv("POSTGRES_HOST"),
				os.Getenv("POSTGRES_PORT"),
				os.Getenv("POSTGRES_DB"))
	}

	rssURL := flag.String("url", "", "RSS feed URL")
	dbParams := flag.String("db", defaultConnection, "Postgres connection string")
	limit := flag.Int("limit", 0, "Max number of entries (0 = all)")
	reset := flag.Bool("reset", false, "Clear table before inserting")
	showDB := flag.Bool("show-db", false, "Show table contents")
	flag.Parse()

	if *rssURL == "" && !*showDB && !*reset {
		log.Fatal("Please specify either --url or --show-db or --reset parameter")
	}

	dbConn, err := db.InitDB(*dbParams)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close(dbConn)

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

		var count int
		for _, item := range items {
			if *limit > 0 && count >= *limit {
				break
			}
			log.Printf("Storing: %s ...", item.Title)

			err := db.SaveItem(dbConn, item.Title, item.Link, item.Description, item.Author, item.Category, item.PubDate, item.GUID)

			if err != nil {
				log.Printf("Insert error: %v", err)
				continue
			}
			count++
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

		var count int
		for _, item := range items {
			fmt.Printf("%d | %s | %s\n", item.ID, item.Title, item.Link)
			count++
		}

		log.Println(SPLIT_LINE)
		log.Println("== Found ", count, "items in the database")
		log.Println(SPLIT_LINE)
	}
}
