-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    description TEXT,
    published_at TIMESTAMP,
    feed_id INTEGER NOT NULL,
    FOREIGN KEY (feed_id)
    REFERENCES feeds(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
    DROP TABLE posts;
-- +goose StatementEnd
