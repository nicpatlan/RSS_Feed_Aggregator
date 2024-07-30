-- name: CreatePost :one
INSERT INTO posts
    (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsByUser :many
SELECT p.id, p.created_at, p.updated_at, p.title, p.url, p.description, p.published_at, p.feed_id
FROM posts p
    INNER JOIN feeds ON p.feed_id = feeds.id
WHERE feeds.user_id = $1
ORDER BY p.published_at DESC
LIMIT $2;