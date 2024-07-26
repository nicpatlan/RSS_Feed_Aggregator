package main

import (
	"net/http"

	"github.com/nicpatlan/RSS_Feed_Aggregator/internal/database"
)

type authenticatedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiConfig *apiConfig) authMiddlewareHandler(handler authenticatedHandler) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		apiKey, err := getRequestApiKey(req)
		if err != nil {
			respondWithError(wr, http.StatusBadRequest, err.Error())
			return
		}
		user, err := apiConfig.DB.GetUserByApiKey(req.Context(), apiKey)
		if err != nil {
			respondWithError(wr, http.StatusBadRequest, err.Error())
			return
		}
		handler(wr, req, user)
	}
}
