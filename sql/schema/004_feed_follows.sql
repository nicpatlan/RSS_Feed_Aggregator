-- +goose Up
CREATE TABLE users_feeds
(
    id UUID PRIMARY KEY,
    feed_id UUID NOT NULL REFERENCES feeds(id),
    user_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE(user_id, feed_id)
);

-- +goose Down
DROP TABLE users_feeds;