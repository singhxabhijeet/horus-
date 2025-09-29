package main

import (
	"log"
	"net/http"

	"github.com/singhxabhijeet/horus/internal" // Use your actual module path
)

func main() {
	// Load configuration from .env file or environment
	cfg := LoadConfig()

	// Connect to the database
	db := internal.NewDB(cfg.DatabaseURL)
	defer db.Close()

	// Create a new ServeMux (router)
	mux := http.NewServeMux()

	// Set up API handlers
	api := internal.NewAPI(db)
	api.RegisterRoutes(mux)

	// Start the background health checker
	internal.StartChecker(db)

	// Start the HTTP server
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
