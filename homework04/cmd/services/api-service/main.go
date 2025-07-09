package main

import (
	"flag"
	"log"
	"rss_fetcher/internal/db"
	"rss_fetcher/internal/services/api"

	"github.com/labstack/echo/v4"
)

const (
	rootPath     = "/api/v1"
	channelsPath = rootPath + "/channels"
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

	api := api.New(dbConn, vectorDB, "http://localhost:11434", "mxbai-embed-large", "llama3")

	e := echo.New()

	e.POST(channelsPath, api.AddChannel)
	e.GET(channelsPath, api.GetChannels)
	e.GET(channelsPath+"/:id", api.GetChannel)
	e.DELETE(channelsPath, api.DeleteChannels)
	e.DELETE(channelsPath+"/:id", api.DeleteChannel)
	e.GET(rootPath+"/news", api.GetAllNews)
	e.DELETE(rootPath+"/news/:id", api.DeleteNews)
	e.GET(rootPath+"/query/:q", api.GetQuery)
	e.GET(rootPath+"/jobs", api.GetJobs)

	log.Println("Starting public service on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}
