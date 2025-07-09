package daemon

import (
	"context"
	"log"
	"rss_fetcher/internal/data"
	"rss_fetcher/internal/db"
	"rss_fetcher/internal/parser"
	"time"

	backend "rss_fetcher/internal/ollama"

	"github.com/jackc/pgx/v5"
)

type NewsDaemon struct {
	dbConn     *pgx.Conn
	vectorDB   *db.PGVector
	ollamaHost string
	embModel   string
	genModel   string
}

func NewNewsDaemon(dbConn *pgx.Conn, vectorDB *db.PGVector, host, embModel, genModel string) *NewsDaemon {
	return &NewsDaemon{dbConn: dbConn, vectorDB: vectorDB, embModel: embModel, genModel: genModel, ollamaHost: host}
}

func (daemon *NewsDaemon) CheckJobs() {
	items, err := db.LoadJobs(db.NewConnectionQuery(daemon.dbConn), " Status = 0")
	if err != nil {
		log.Printf("Failed to load jobs: %v", err)
		return
	}

	for _, item := range items {
		if err := daemon.processJob(item); err != nil {
			log.Printf("Error processing job for %s news: %v", item.Link, err)
		}
	}
}

func (daemon *NewsDaemon) saveChunk(newsID int, chunk string) error {
	embeddingBackend := backend.NewOllamaBackend(daemon.ollamaHost, daemon.embModel, time.Duration(20*time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	embedding, err := embeddingBackend.Embed(ctx, chunk, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		log.Printf("Failed to generate embedding for chunk: %s, error: %v", chunk, err)
		return err
	}
	err = daemon.vectorDB.InsertDocument(ctx, newsID, chunk, embedding)
	if err != nil {
		log.Fatalf("Error inserting data into vector database: %v", err)
		return err
	}
	return nil
}

func (daemon *NewsDaemon) processJob(job data.NewsJob) error {
	err := db.UpdateJobStatus(db.NewConnectionQuery(daemon.dbConn), job.ID, 1)
	if err != nil {
		log.Printf("Failed to update job status for job %d: %v", job.ID, err)
		return err
	}

	defer db.UpdateJobStatus(db.NewConnectionQuery(daemon.dbConn), job.ID, 2)

	linkItems, err := db.LoadNewsByLink(db.NewConnectionQuery(daemon.dbConn), job.Link)
	if err != nil || len(linkItems) == 0 {
		log.Printf("Failed to load news for link %s: %v", job.Link, err)
	}

	news := linkItems[0]

	chunks, err := parser.ExtractArticle(job.Link, 400)
	if err != nil {
		log.Printf("Failed to parse article %s: %v", job.Link, err)
		return err
	}

	for _, chunk := range chunks {
		daemon.saveChunk(news.ID, chunk)
		if err != nil {
			log.Printf("Failed to save chunk %s for link %s: %v", chunk, news.Link, err)
		}
	}
	log.Printf("Processed news.link %s with %d chunks", news.Link, len(chunks))
	return nil
}
