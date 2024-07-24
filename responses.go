package main

import (
	"encoding/json"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func respondWithJSON(wr http.ResponseWriter, statusCode int, payload interface{}) {
	wr.Header().Set("Content-Type", "application/json")
	validResponse, err := json.Marshal(payload)
	if err != nil {
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}
	wr.WriteHeader(statusCode)
	wr.Write([]byte(validResponse))
}

func respondWithError(wr http.ResponseWriter, statusCode int, msg string) {
	errResponse := errorResponse{
		Error: msg,
	}
	respondWithJSON(wr, statusCode, errResponse)
}
