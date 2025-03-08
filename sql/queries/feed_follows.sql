-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (INSERT INTO feed_follows(created_at, updated_at, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *)

SELECT 
    inserted_feed_follow.*,
    users.name as user_name,
    feeds.name as feed_name
FROM inserted_feed_follow
INNER JOIN users ON users.id = inserted_feed_follow.user_id
INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id;

-- name: GetFeedFollowsForUser :many
SELECT 
    users.name as username,
    feeds.name as feed_follows
FROM feeds
INNER JOIN feed_follows ON feeds.id = feed_follows.feed_id
INNER JOIN users ON feed_follows.user_id = users.id
WHERE users.name = $1;

-- name: RemoveFeedFollow :exec
DELETE FROM feed_follows
WHERE user_id = (
    SELECT id FROM users
    WHERE users.name = $1
) AND
feed_id = (
    SELECT id FROM feeds
    WHERE url = $2
);
