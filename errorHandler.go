package main

import (
	"net/http"
)

type errorHandler struct{}

func (errorHandler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	respondWithError(wr, http.StatusInternalServerError, "error")
}
