package main

import (
	"net/http"
)

func (apiConfig *apiConfig) getFeedsHandler(wr http.ResponseWriter, req *http.Request) {
	type feedsResponse struct {
		Feeds []Feed `json:"feeds"`
	}
	dbFeeds, err := apiConfig.DB.GetFeeds(req.Context())
	if err != nil {
		respondWithError(wr, http.StatusInternalServerError, err.Error())
		return
	}
	feeds := feedsResponse{
		Feeds: make([]Feed, 0),
	}
	for _, dbFeed := range dbFeeds {
		feeds.Feeds = append(feeds.Feeds, convertDatabaseFeedToFeed(dbFeed))
	}
	respondWithJSON(wr, http.StatusOK, feeds)
}
