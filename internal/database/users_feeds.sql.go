// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users_feeds.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const followFeed = `-- name: FollowFeed :one
INSERT INTO users_feeds
    (id, feed_id, user_id, created_at, updated_at)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING id, feed_id, user_id, created_at, updated_at
`

type FollowFeedParams struct {
	ID        uuid.UUID
	FeedID    uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) FollowFeed(ctx context.Context, arg FollowFeedParams) (UsersFeed, error) {
	row := q.db.QueryRowContext(ctx, followFeed,
		arg.ID,
		arg.FeedID,
		arg.UserID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i UsersFeed
	err := row.Scan(
		&i.ID,
		&i.FeedID,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const unfollowFeed = `-- name: UnfollowFeed :exec
DELETE FROM users_feeds
WHERE id = $1
`

func (q *Queries) UnfollowFeed(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, unfollowFeed, id)
	return err
}
