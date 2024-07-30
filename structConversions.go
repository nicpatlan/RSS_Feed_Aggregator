package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/nicpatlan/RSS_Feed_Aggregator/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

type Feed struct {
	ID            uuid.UUID  `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Name          string     `json:"name"`
	Url           string     `json:"url"`
	UserID        uuid.UUID  `json:"user_id"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
}

type UsersFeed struct {
	ID        uuid.UUID `json:"id"`
	FeedID    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func convertDatabaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}
}

func convertDatabaseFeedToFeed(feed database.Feed) Feed {
	var last_fetched_time *time.Time
	if feed.LastFetchedAt.Valid {
		last_fetched_time = &feed.LastFetchedAt.Time
	} else {
		last_fetched_time = nil
	}
	return Feed{
		ID:            feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		Name:          feed.Name,
		Url:           feed.Url,
		UserID:        feed.UserID,
		LastFetchedAt: last_fetched_time,
	}
}

func convertDatabaseFeedsToArray(dbFeeds []database.Feed) []Feed {
	feeds := make([]Feed, len(dbFeeds))
	for idx, feed := range dbFeeds {
		feeds[idx] = convertDatabaseFeedToFeed(feed)
	}
	return feeds
}

func convertDatabaseUsersFeedToUsersFeed(usersFeed database.UsersFeed) UsersFeed {
	return UsersFeed{
		ID:        usersFeed.ID,
		FeedID:    usersFeed.FeedID,
		UserID:    usersFeed.UserID,
		CreatedAt: usersFeed.CreatedAt,
		UpdatedAt: usersFeed.UpdatedAt,
	}
}

func convertDatabaseUsersFeedToArray(dbUsersFeeds []database.UsersFeed) []UsersFeed {
	usersFeeds := make([]UsersFeed, len(dbUsersFeeds))
	for idx, usersFeed := range dbUsersFeeds {
		usersFeeds[idx] = convertDatabaseUsersFeedToUsersFeed(usersFeed)
	}
	return usersFeeds
}

func convertDatabasePostToPost(post database.Post) Post {
	return Post{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Url:         post.Url,
		Description: post.Description,
		PublishedAt: post.PublishedAt,
		FeedID:      post.FeedID,
	}
}

func convertDatabasePostsToArray(posts []database.Post) []Post {
	userPosts := make([]Post, len(posts))
	for idx, post := range posts {
		userPosts[idx] = convertDatabasePostToPost(post)
	}
	return userPosts
}
