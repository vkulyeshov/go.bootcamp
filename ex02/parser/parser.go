package parser

import (
	"github.com/mmcdole/gofeed"
)

type RSSItem struct {
	Title string
	Link  string
}

func FetchRSS(url string) ([]RSSItem, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}

	var items []RSSItem
	for _, entry := range feed.Items {
		items = append(items, RSSItem{
			Title: entry.Title,
			Link:  entry.Link,
		})
	}

	return items, nil
}
