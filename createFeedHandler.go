package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nicpatlan/RSS_Feed_Aggregator/internal/database"
)

func (apiConfig *apiConfig) createFeedHandler(wr http.ResponseWriter, req *http.Request, user database.User) {
	// get feed request parameters
	type feedRequest struct {
		Name string
		URL  string
	}
	decoder := json.NewDecoder(req.Body)
	feedReq := feedRequest{}
	err := decoder.Decode(&feedReq)
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}

	// add feed to database
	createParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      feedReq.Name,
		Url:       feedReq.URL,
		UserID:    user.ID,
	}
	feed, err := apiConfig.DB.CreateFeed(req.Context(), createParams)
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}

	// set user to follow created feed automatically
	followParams := database.FollowFeedParams{
		ID:        uuid.New(),
		FeedID:    feed.ID,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	usersFeed, err := apiConfig.DB.FollowFeed(req.Context(), followParams)
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}

	// set up response object
	type response struct {
		Feed       Feed      `json:"feed"`
		FeedFollow UsersFeed `json:"feed_follow"`
	}
	createResponse := response{
		Feed:       convertDatabaseFeedToFeed(feed),
		FeedFollow: convertDatabaseUsersFeedToUsersFeed(usersFeed),
	}
	respondWithJSON(wr, http.StatusCreated, createResponse)
}
