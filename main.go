package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// get environment variables
	godotenv.Load()
	port := os.Getenv("PORT")

	// create server multiplexer and http server
	serveMux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	// endpoint patterns
	const healthzPattern = "GET /v1/healthz"
	const errorPattern = "GET /v1/err"

	// add handlers
	serveMux.Handle(healthzPattern, healthHandler{})
	serveMux.Handle(errorPattern, errorHandler{})

	// start server
	fmt.Printf("Running server and listening on port: %s", port)
	log.Fatal(server.ListenAndServe())
}
