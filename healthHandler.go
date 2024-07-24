package main

import (
	"net/http"
)

type healthHandler struct{}

func (healthHandler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	response := healthResponse{
		Status: "ok",
	}
	respondWithJSON(wr, http.StatusOK, response)
}
