# RSS Rag system

A simple tool what fetching last neww based on RSS channels
Then it makes analyse and store relevant information to Vector DB
It allows to produce LLM answers based on last news

## How to use

Firstly run up services
```bash
docker compose up -d --build
```

Then run cli app can be run:

```bash
go -C app run ./cmd/cli/ <options>
```

## Features

- Add RSS channel to db (`--url`)
- Clear the database (`--reset`)
- View channels from the database (`--show-channels`)
- View news from the database (`--show-news`)
- View jobs from the database (`--show-jobs`)
- Chat with LLM (`--query="query text"`)

---



# List of available RSS resources as example for testing purpose
https://about.fb.com/wp-content/uploads/2016/05/rss-urls-1.pdf