package main

import (
	"flag"
	"log"
	"rss_fetcher/internal/db"
	"rss_fetcher/internal/services/daemon"
	"time"
)

const (
	defaultConnection       = "postgres://rss:rss@localhost:5432/rss"
	ollamaDefaultConnection = "http://localhost:11434"
	embDefaultModel         = "mxbai-embed-large"
	genDefaultModel         = "llama3"
)

func main() {
	dbParams := flag.String("db", defaultConnection, "Postgres connection string")
	ollamaConnection := flag.String("ollama", ollamaDefaultConnection, "Postgres connection string")
	embModel := flag.String("emb", embDefaultModel, "Embedding model")
	genModel := flag.String("gen", genDefaultModel, "Generative model")

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

	daemon := daemon.NewNewsDaemon(dbConn, vectorDB, *ollamaConnection, *embModel, *genModel)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("Starting news daemon")
	for {
		daemon.CheckJobs()
		<-ticker.C
	}
}
