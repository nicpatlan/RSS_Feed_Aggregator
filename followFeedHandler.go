package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nicpatlan/RSS_Feed_Aggregator/internal/database"
)

func (apiConfig *apiConfig) followFeedHandler(wr http.ResponseWriter, req *http.Request, user database.User) {
	// get the body of request
	type followRequest struct {
		FeedID string `json:"feed_id"`
	}
	decoder := json.NewDecoder(req.Body)
	followReq := followRequest{}
	err := decoder.Decode(&followReq)
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}

	// add follow to the users_feeds database
	feedID, err := uuid.Parse(followReq.FeedID)
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}
	params := database.FollowFeedParams{
		ID:        uuid.New(),
		FeedID:    feedID,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	usersFeed, err := apiConfig.DB.FollowFeed(req.Context(), params)
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(wr, http.StatusCreated, convertDatabaseUsersFeedToUsersFeed(usersFeed))
}
