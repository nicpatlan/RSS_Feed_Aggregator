package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nicpatlan/RSS_Feed_Aggregator/internal/database"
)

type userRequest struct {
	Name string `json:"name"`
}

func (apiConfig *apiConfig) createUserHandler(wr http.ResponseWriter, req *http.Request) {
	// get the create user request body
	decoder := json.NewDecoder(req.Body)
	userReq := userRequest{}
	decoder.Decode(&userReq)

	// generate a new UUID
	newUUID := uuid.New()

	// create params and insert user in database
	newUser := database.CreateUserParams{
		ID:        newUUID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      userReq.Name,
	}
	user, err := apiConfig.DB.CreateUser(req.Context(), newUser)
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(wr, http.StatusCreated, convertDatabaseUserToUser(user))
}
