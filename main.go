package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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

	// add handlers
	serveMux.Handle(healthzPattern, healthHandler{})
	serveMux.Handle(errorPattern, errorHandler{})
	serveMux.HandleFunc(createUserPattern, apiConfig.createUserHandler)

	// start server
	fmt.Printf("Running server and listening on port: %s", port)
	log.Fatal(server.ListenAndServe())
}
