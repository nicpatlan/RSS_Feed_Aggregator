package main

import (
	"net/http"
)

func (apiConfig *apiConfig) getFeedsHandler(wr http.ResponseWriter, req *http.Request) {
	dbFeeds, err := apiConfig.DB.GetFeeds(req.Context())
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(wr, http.StatusOK, convertDatabaseFeedsToArray(dbFeeds))
}
