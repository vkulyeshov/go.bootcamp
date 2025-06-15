# RSS Fetcher

A simple CLI tool written in Go for fetching RSS feeds, storing them in an PostgreSQL database, and viewing saved entries.

## Features

- Fetch RSS feed from a URL (`--url`)
- Limit the number of items fetched from RSS resource (`--limit`)
- Store feed items to PostgreSQL database (`--db`) default `postgres://rss:rss@localhost:5432/rss`
- Clear the database (`--reset`)
- View entries from the database (`--show-db`)

---

## Build

```bash
go -C app build -o rss-fetcher
```

## Run

```bash
go -C app run . <options>
```

# List of available RSS resources as example for testing purpose
https://about.fb.com/wp-content/uploads/2016/05/rss-urls-1.pdf