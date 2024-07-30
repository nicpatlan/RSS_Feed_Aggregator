package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nicpatlan/RSS_Feed_Aggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// get environment variables
	godotenv.Load()
	port := os.Getenv("PORT")
	dbURL := os.Getenv("CONN")

	// open connection to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("could not load database")
	}
	dbQueries := database.New(db)
	apiConfig := apiConfig{
		DB: dbQueries,
	}

	// create server multiplexer and http server
	serveMux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	// endpoint patterns
	const healthzPattern = "GET /v1/healthz"
	const errorPattern = "GET /v1/err"
	const createUserPattern = "POST /v1/users"
	const getUserApiPattern = "GET /v1/users"
	const createFeedPattern = "POST /v1/feeds"
	const getFeedsPattern = "GET /v1/feeds"
	const followFeedPattern = "POST /v1/feed_follows"
	const unfollowFeedPattern = "DELETE /v1/feed_follows/{feedFollowID}"
	const getUserFeedsPattern = "GET /v1/feed_follows"

	// add handlers
	serveMux.Handle(healthzPattern, healthHandler{})
	serveMux.Handle(errorPattern, errorHandler{})
	serveMux.HandleFunc(createUserPattern, apiConfig.createUserHandler)
	serveMux.HandleFunc(getUserApiPattern, apiConfig.authMiddlewareHandler(apiConfig.getUserWithApiHandler))
	serveMux.HandleFunc(createFeedPattern, apiConfig.authMiddlewareHandler(apiConfig.createFeedHandler))
	serveMux.HandleFunc(getFeedsPattern, apiConfig.getFeedsHandler)
	serveMux.HandleFunc(followFeedPattern, apiConfig.authMiddlewareHandler(apiConfig.followFeedHandler))
	serveMux.HandleFunc(unfollowFeedPattern, apiConfig.authMiddlewareHandler(apiConfig.unfollowFeedHandler))
	serveMux.HandleFunc(getUserFeedsPattern, apiConfig.authMiddlewareHandler(apiConfig.getUserFeedsHandler))

	// start goroutine to fetch a batch of 2 feeds every 60 seconds
	const batchSize = 2
	const fetchInterval = time.Minute
	go apiConfig.fetchFeedBatch(batchSize, fetchInterval)

	// start server
	fmt.Printf("Running server and listening on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
