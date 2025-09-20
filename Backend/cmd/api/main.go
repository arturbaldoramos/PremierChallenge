package main

import (
	"log"
	"net/http"

	"example.com/m/v2/internal/config"
	"example.com/m/v2/internal/database"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting Premiere Challenge Backend...")

	// Initialize database connection
	db := config.InitDatabase()

	// Run automatic migrations
	err := database.RunMigrations(db)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Setup HTTP router
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "running", "message": "Premiere Challenge API"}`))
	}).Methods("GET")

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}