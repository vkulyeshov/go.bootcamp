package parser

import (
	"strings"
	"time"

	prose "github.com/jdkato/prose/v2"

	readability "github.com/go-shiori/go-readability"
)

func countTokens(text string) int {
	return len(strings.Fields(text))
}

func chunkTextByTokens(text string, maxTokens int) ([]string, error) {
	doc, err := prose.NewDocument(text)
	if err != nil {
		return nil, err
	}

	var chunks []string
	var currentChunk []string
	var currentLen int

	for _, sent := range doc.Sentences() {
		sentText := sent.Text
		sentTokens := countTokens(sentText)

		if currentLen+sentTokens > maxTokens {
			chunks = append(chunks, strings.Join(currentChunk, " "))
			currentChunk = []string{}
			currentLen = 0
		}

		currentChunk = append(currentChunk, sentText)
		currentLen += sentTokens
	}

	if len(currentChunk) > 0 {
		chunks = append(chunks, strings.Join(currentChunk, " "))
	}

	return chunks, nil
}

func ExtractArticle(url string, chunkSize int) ([]string, error) {
	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		return nil, err
	}

	chunks, err := chunkTextByTokens(article.TextContent, chunkSize)
	if err != nil {
		return nil, err
	}

	return chunks, nil
}
