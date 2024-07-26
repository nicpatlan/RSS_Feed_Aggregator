package main

import (
	"net/http"

	"github.com/nicpatlan/RSS_Feed_Aggregator/internal/database"
)

func (apiConfig *apiConfig) getUserFeedsHandler(wr http.ResponseWriter, req *http.Request, user database.User) {
	userFeeds, err := apiConfig.DB.GetUserFeeds(req.Context(), user.ID)
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(wr, http.StatusOK, convertDatabaseUsersFeedToArray(userFeeds))
}
