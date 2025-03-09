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

-- name: GetFeedByUrl :one
SELECT id FROM feeds
WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $1, updated_at = $2
WHERE id = $3;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;
