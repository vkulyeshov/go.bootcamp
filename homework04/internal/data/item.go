package data

import (
	"time"
)

type ChannelNews struct {
	ID          int
	Title       string
	Link        string
	Description string
	Author      string
	Category    string
	PubDate     time.Time
	GUID        string
}

type Channel struct {
	ID          int
	Link        string        // URL of the channel
	Title       string        // Title of the channel
	Description string        // Description of the channel
	RSSLink     string        // Link inside the RSS feed
	UpdatedAt   time.Time     // Last updated time of the channel
	Items       []ChannelNews // List of items in the channel
}

type NewsJob struct {
	ID        int
	Link      string
	Status    int // 0 - pending, 1 - in progress, 2 - completed
	CreatedAt time.Time
	UpdatedAt time.Time
}
