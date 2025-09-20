package main

import (
	"log"
	"net/http"

	"example.com/m/v2/internal/config"
	"example.com/m/v2/internal/database"
	"example.com/m/v2/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting Premiere Challenge Backend...")

	// Initialize database connection
	db := config.InitDatabase()

	// Reset database (temporary for fixing migration issues)
	err := database.ResetDatabase(db)
	if err != nil {
		log.Fatalf("Failed to reset database: %v", err)
	}

	// Run automatic migrations
	err = database.RunMigrations(db)
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

	// Initialize handlers
	uploadHandler := handlers.NewUploadHandler(db)

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "running", "message": "Premiere Challenge API"}`))
	}).Methods("GET")

	// Upload routes for each domain
	api.HandleFunc("/upload/estados", uploadHandler.UploadEstados).Methods("POST")
	api.HandleFunc("/upload/municipios", uploadHandler.UploadMunicipios).Methods("POST")
	api.HandleFunc("/upload/hospitais", uploadHandler.UploadHospitais).Methods("POST")
	api.HandleFunc("/upload/pacientes", uploadHandler.UploadPacientes).Methods("POST")
	api.HandleFunc("/upload/medicos", uploadHandler.UploadMedicos).Methods("POST")
	api.HandleFunc("/upload/cid10", uploadHandler.UploadCID10).Methods("POST")

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}