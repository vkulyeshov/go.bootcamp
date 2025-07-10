package main

import (
	"flag"
	"fmt"
	"log"
	"rss_fetcher/internal/db"
	"rss_fetcher/internal/services/api"

	"github.com/labstack/echo/v4"
)

const (
	defaultConnection       = "postgres://rss:rss@localhost:5432/rss"
	ollamaDefaultConnection = "http://localhost:11434"
	embDefaultModel         = "mxbai-embed-large"
	genDefaultModel         = "llama3"
	defaultApiPort          = 8080

	rootPath     = "/api/v1"
	channelsPath = rootPath + "/channels"
	newsPath     = rootPath + "/news"
	jobsPath     = rootPath + "/jobs"
	queryPath    = rootPath + "/query"
)

func main() {

	dbParams := flag.String("db", defaultConnection, "Postgres connection string")
	ollamaConnection := flag.String("ollama", ollamaDefaultConnection, "Postgres connection string")
	embModel := flag.String("emb", embDefaultModel, "Embedding model")
	genModel := flag.String("gen", genDefaultModel, "Generative model")
	apiPort := flag.Int("port", defaultApiPort, "REST API port")

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

	api := api.New(dbConn, vectorDB, *ollamaConnection, *embModel, *genModel)

	e := echo.New()

	e.POST(channelsPath, api.AddChannel)
	e.GET(channelsPath, api.GetChannels)
	e.GET(channelsPath+"/:id", api.GetChannel)
	e.DELETE(channelsPath, api.DeleteChannels)
	e.DELETE(channelsPath+"/:id", api.DeleteChannel)
	e.GET(newsPath, api.GetAllNews)
	e.DELETE(newsPath+"/:id", api.DeleteNews)
	e.GET(queryPath+"/:q", api.GetQuery)
	e.GET(jobsPath, api.GetJobs)

	// graceful exit from service
	//quitChannel := make(chan os.Signal, 1)
	//signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutines
	//go func() {
	log.Printf("Starting REST API service on %d", *apiPort)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *apiPort)))
	//}()

	// Block until a signal is received
	//<-quitChannel
	//log.Println("Shutting down REST API service...")

	// Create a context with a timeout for graceful shutdown
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

	// Attempt graceful shutdown
	//if err := e.Shutdown(ctx); err != nil {
	//	log.Fatalf("REST API service shutdown error: %v", err)
	//}

	//log.Println("REST API service is stopped.")
}
