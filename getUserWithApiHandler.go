package main

import (
	"net/http"
	"strings"
)

func (apiConfig *apiConfig) getUserWithApiHandler(wr http.ResponseWriter, req *http.Request) {
	authStr := req.Header.Get("Authorization")
	_, apiKey, ok := strings.Cut(authStr, "ApiKey ")
	if !ok {
		respondWithError(wr, http.StatusBadRequest, "invalid header")
		return
	}
	user, err := apiConfig.DB.GetUserByApiKey(req.Context(), apiKey)
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}
	userResponse := convertDatabaseUserToUser(user)
	respondWithJSON(wr, http.StatusOK, userResponse)
}
