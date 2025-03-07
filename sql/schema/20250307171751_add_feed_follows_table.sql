-- +goose Up
-- +goose StatementBegin
CREATE TABLE feed_follows(
        id SERIAL PRIMARY KEY,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        user_id UUID NOT NULL,
        feed_id INTEGER NOT NULL,
        FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
        FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
        CONSTRAINT unique_user_feed UNIQUE (user_id, feed_id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feed_follows;
-- +goose StatementEnd
