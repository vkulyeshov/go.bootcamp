package main

import (
	"flag"
	"log"
	"rss_fetcher/internal/db"
	"rss_fetcher/internal/services/daemon"
	"time"
)

const (
	defaultConnection = "postgres://rss:rss@localhost:5432/rss"
)

func main() {
	dbParams := flag.String("db", defaultConnection, "Postgres connection string")
	flag.Parse()
	dbConn, err := db.InitDB(*dbParams)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close(dbConn)

	daemon := daemon.NewChannelDaemon(dbConn)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("Starting channel daemon")
	for {
		daemon.CheckFeeds()
		<-ticker.C
	}
}
