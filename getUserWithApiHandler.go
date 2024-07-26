package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/nicpatlan/RSS_Feed_Aggregator/internal/database"
)

func getRequestApiKey(req *http.Request) (string, error) {
	authStr := req.Header.Get("Authorization")
	_, apiKey, ok := strings.Cut(authStr, "ApiKey ")
	if !ok {
		return "", errors.New("invalid header")
	}
	return apiKey, nil
}

func (apiConfig *apiConfig) getUserWithApiHandler(wr http.ResponseWriter, req *http.Request, user database.User) {
	respondWithJSON(wr, http.StatusOK, convertDatabaseUserToUser(user))
}
