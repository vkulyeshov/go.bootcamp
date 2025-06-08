package parser

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"rss_fetcher/data"
)

type xmlRSS struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel xmlChannel `xml:"channel"`
}
type xmlChannel struct {
	Title         string       `xml:"title"`
	Link          string       `xml:"link"`
	Description   string       `xml:"description"`
	Language      string       `xml:"language"`
	Copyright     string       `xml:"copyright"`
	PubDate       string       `xml:"pubDate"`
	LastBuildDate string       `xml:"lastBuildDate"`
	Generator     string       `xml:"generator"`
	Items         []xmlRSSItem `xml:"item"`
}

type xmlRSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Author      string `xml:"author"`
	Category    string `xml:"category"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

func FetchRSS(url string) ([]data.RSSItem, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch RSS feed: non-200 response")
	}

	// Read all as first implementation
	// better to use a streaming parser for large feeds and pass it to db writer as some interface
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var xmlData xmlRSS
	if err := xml.Unmarshal(body, &xmlData); err != nil {
		return nil, err
	}

	// Convert xmlRSSItem to data.RSSItem
	items := make([]data.RSSItem, len(xmlData.Channel.Items))

	for i, item := range xmlData.Channel.Items {
		items[i] = data.RSSItem{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Author:      item.Author,
			Category:    item.Category,
			PubDate:     item.PubDate,
			GUID:        item.GUID,
		}
	}
	return items, nil
}
