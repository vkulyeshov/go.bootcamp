package api

import (
	"context"
	"log"
	"net/http"
	"rss_fetcher/internal/db"
	"strconv"
	"time"

	backend "rss_fetcher/internal/ollama"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
)

type addChannelRequest struct {
	Link string `json:"link"`
}

type API struct {
	dbConn     *pgx.Conn
	vectorDB   *db.PGVector
	ollamaHost string
	embModel   string
	genModel   string
}

func New(dbConn *pgx.Conn, vectorDB *db.PGVector, host, embModel, genModel string) *API {
	return &API{dbConn: dbConn, vectorDB: vectorDB, embModel: embModel, genModel: genModel, ollamaHost: host}
}

func (api *API) AddChannel(c echo.Context) error {
	var channel addChannelRequest
	if err := c.Bind(&channel); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	if channel.Link == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Link is required"})
	}
	if err := db.AddChannel(api.dbConn, channel.Link); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				log.Printf("Channel with link: %v already added to the channels list", channel.Link)
				return c.JSON(http.StatusConflict, map[string]string{"error": "Channel is already exists"})
			}
			log.Printf("Postgres error: %s\n", pgErr.Message)
		}
		log.Printf("Failed to add channel: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add channel"})
	}
	return c.JSON(http.StatusCreated, map[string]string{"status": "ok"})
}

func (api *API) GetChannels(c echo.Context) error {
	channels, err := db.LoadChannels(api.dbConn)
	if err != nil {
		log.Printf("Error loading channels from database: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load channels"})
	}

	return c.JSON(http.StatusOK, channels)
}

func (api *API) GetChannel(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Wrong id: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid channel ID	"})
	}
	result, err := db.LoadChannel(api.dbConn, id)
	if err != nil {
		log.Printf("Error loading channel from database: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load channel"})
	}
	return c.JSON(http.StatusOK, result)
}

func (api *API) DeleteChannels(c echo.Context) error {
	if err := db.DeleteChannels(api.dbConn); err != nil {
		log.Printf("Error deleting channels: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete channels"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "All channels deleted"})
}

func (api *API) DeleteChannel(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Wrong channel id: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid channel ID"})
	}
	if err := db.DeleteChannel(api.dbConn, id); err != nil {
		log.Printf("Error deleting channel: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to	delete channel"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Channel deleted"})
}

func (api *API) DeleteNews(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Wrong news id: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid news ID"})
	}
	if err := db.DeleteNews(api.dbConn, id); err != nil {
		log.Printf("Error in deleting news: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to	delete news"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "News deleted"})
}

func (api *API) GetAllNews(c echo.Context) error {
	result, err := db.LoadNews(db.NewConnectionQuery(api.dbConn))
	if err != nil {
		log.Printf("Error loading news from database: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load news"})
	}

	return c.JSON(http.StatusOK, result)
}

func (api *API) GetJobs(c echo.Context) error {
	items, err := db.LoadJobs(db.NewConnectionQuery(api.dbConn), "")
	if err != nil {
		log.Printf("Error loading channels from database: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load jobs"})
	}
	return c.JSON(http.StatusOK, items)
}

func (api *API) GetQuery(c echo.Context) error {
	q := c.Param("q")

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	embeddingBackend := backend.NewOllamaBackend(api.ollamaHost, api.embModel, time.Duration(60*time.Second))
	generationBackend := backend.NewOllamaBackend(api.ollamaHost, api.genModel, time.Duration(60*time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Embed the query using the specified embedding backend
	queryEmbedding, err := embeddingBackend.Embed(ctx, q, headers)
	if err != nil {
		log.Printf("Error generating query embedding: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error generating query embedding"})
	}
	log.Println("Vector embeddings generated")

	// Retrieve relevant documents for the query embedding
	retrievedDocs, err := api.vectorDB.QueryRelevantDocuments(ctx, queryEmbedding, "ollama")
	if err != nil {
		log.Printf("Error retrieving relevant documents: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error retrieving relevant documents"})
	}

	// Log the retrieved documents to see if they include the inserted content
	for _, doc := range retrievedDocs {
		log.Printf("Retrieved Document: %v", doc)
	}

	// Augment the query with retrieved context
	augmentedQuery := db.CombineQueryWithContext(q, retrievedDocs)

	log.Println("READY FOR GENERATING PROMPT")

	prompt := backend.NewPrompt().
		AddMessage("system", "You — are AI assistent. Use the provided context to answer the user question").
		AddMessage("user", augmentedQuery).
		SetParameters(backend.Parameters{
			MaxTokens:   150, // Supported by LLaMa
			Temperature: 0.7, // Supported by LLaMa
			TopP:        0.9, // Supported by LLaMa
		})

	log.Println("PROMPT PREAPARED")
	// Generate response with the specified generation backend
	response, err := generationBackend.Generate(ctx, prompt)
	if err != nil {
		log.Printf("Error generating response: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error generating response"})
	}

	return c.JSON(http.StatusOK,
		map[string]string{
			"query":    q,
			"response": response,
		})
}
