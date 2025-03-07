-- name: CreateFeed :one
INSERT INTO feeds (created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: ResetFeed :exec
DELETE FROM feeds;

-- name: GetFeeds :many
SELECT feeds.name, feeds.url, users.name as creator FROM feeds
LEFT JOIN users ON users.id = feeds.user_id;

