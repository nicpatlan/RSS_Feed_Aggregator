-- name: FollowFeed :one
INSERT INTO users_feeds
    (id, feed_id, user_id, created_at, updated_at)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UnfollowFeed :exec
DELETE FROM users_feeds
WHERE id = $1;