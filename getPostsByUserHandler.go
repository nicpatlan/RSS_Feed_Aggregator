package main

import (
	"net/http"
	"strconv"

	"github.com/nicpatlan/RSS_Feed_Aggregator/internal/database"
)

func (apiConfig *apiConfig) getPostsByUserHandler(wr http.ResponseWriter, req *http.Request, user database.User) {
	// get query value and convert to int or set to default of 5
	limitStr := req.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "5"
	}
	queryLimit, err := strconv.Atoi(limitStr)
	if err != nil {
		respondWithError(wr, http.StatusBadRequest, err.Error())
		return
	}

	// set parameters for post retrieval from database
	postsParams := database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(queryLimit),
	}
	dbResp, err := apiConfig.DB.GetPostsByUser(req.Context(), postsParams)
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(wr, http.StatusOK, convertDatabasePostsToArray(dbResp))
}
