-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS vector;
CREATE TABLE IF NOT EXISTS channels (
    channel_id SERIAL PRIMARY KEY,
    link TEXT NOT NULL UNIQUE,
    title TEXT DEFAULT '',
    description TEXT DEFAULT '',
    rss_link TEXT DEFAULT '',
    last_updated timestamp NOT NULL DEFAULT '1970-01-01 00:00:00'
);
CREATE TABLE IF NOT EXISTS channel_news (
    news_id SERIAL PRIMARY KEY,
    channel_id INTEGER NOT NULL REFERENCES channels(channel_id) ON DELETE CASCADE,
    link TEXT NOT NULL UNIQUE,
    title TEXT DEFAULT '',
    description TEXT DEFAULT '',
    author TEXT  DEFAULT '',
    category TEXT  DEFAULT '',
    pub_date timestamp,	
    guid TEXT NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS news_embeddings (
    embedding_id SERIAL PRIMARY KEY,
    news_id INTEGER NOT NULL REFERENCES channel_news(news_id) ON DELETE CASCADE,
    embedding VECTOR(1024),
    metadata JSONB
);
CREATE TABLE IF NOT EXISTS news_jobs (
    job_id SERIAL PRIMARY KEY,
    link TEXT NOT NULL,
    status INTEGER NOT NULL DEFAULT 0, -- 0: pending, 1: in progress, 2: completed
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT '1970-01-01 00:00:00'
);

CREATE INDEX channel_news_pub_date_idx ON channel_news(pub_date); 
CREATE INDEX news_embeddings_idx ON news_embeddings
USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS channel_news_embeddings_idx;
DROP INDEX IF EXISTS channel_news_pub_date_idx;
DROP TABLE channel_jobs;
DROP TABLE news_embeddings;
DROP TABLE channel_items;
DROP TABLE channels;
DROP EXTENSION IF EXISTS vector;
-- +goose StatementEnd