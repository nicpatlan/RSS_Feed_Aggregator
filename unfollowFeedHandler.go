package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/nicpatlan/RSS_Feed_Aggregator/internal/database"
)

func (apiConfig *apiConfig) unfollowFeedHandler(wr http.ResponseWriter, req *http.Request, _ database.User) {
	feedStr := req.PathValue("feedFollowID")
	if feedStr == "" {
		respondWithError(wr, http.StatusBadRequest, "invalid feed ID")
		return
	}
	feedID, err := uuid.Parse(feedStr)
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}
	err = apiConfig.DB.UnfollowFeed(req.Context(), feedID)
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(wr, http.StatusOK, "")
}
