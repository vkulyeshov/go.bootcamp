package main

import (
	"flag"
	"log"
	"rss_fetcher/internal/db"
	"rss_fetcher/internal/services/daemon"
	"time"
)

func main() {
	defaultConnection := "postgres://rss:rss@localhost:5432/rss"
	dbParams := flag.String("db", defaultConnection, "Postgres connection string")
	flag.Parse()
	dbConn, err := db.InitDB(*dbParams)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close(dbConn)

	vectorDB, err := db.NewPGVector(*dbParams)
	if err != nil {
		log.Fatalf("Error initializing vector database: %v", err)
	}
	defer vectorDB.Close()

	daemon := daemon.NewNewsDaemon(dbConn, vectorDB, "http://localhost:11434", "mxbai-embed-large", "llama3")

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("Starting news daemon")
	for {
		daemon.CheckJobs()
		<-ticker.C
	}
}
