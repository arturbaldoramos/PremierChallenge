package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"example.com/m/v2/internal/config"
	"example.com/m/v2/internal/database"
	"example.com/m/v2/internal/handlers"
	"example.com/m/v2/internal/services"
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

	// CORS middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Handle preflight OPTIONS request
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Initialize services
	redisService := services.NewRedisService("localhost:6379", "", 0)
	dataService := services.NewDataService(db)
	workerService := services.NewWorkerService(redisService, dataService)
	monitorService := services.NewMonitorService(redisService)

	// Initialize handlers
	uploadHandler := handlers.NewUploadHandler(db)
	wsHandler := handlers.NewWebSocketHandler(redisService, dataService)

	// Start background services
	workerConfig := services.WorkerConfig{
		MaxWorkers: map[string]int{
			"estados":    2,
			"municipios": 4,
			"medicos":    3,
			"hospitais":  2,
			"pacientes":  3,
			"cid10":      2,
		},
		BatchSize: map[string]int{
			"estados":    50,
			"municipios": 200,
			"medicos":    300,
			"hospitais":  150,
			"pacientes":  250,
			"cid10":      400,
		},
	}

	workerService.StartWorkers(workerConfig)
	monitorService.Start()

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Shutting down services...")
		workerService.Stop()
		monitorService.Stop()
		os.Exit(0)
	}()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "running", "message": "Premiere Challenge API"}`))
	}).Methods("GET")

	// WebSocket endpoint
	router.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// Upload routes for each domain (mantidos para compatibilidade)
	api.HandleFunc("/upload/estados", uploadHandler.UploadEstados).Methods("POST")
	api.HandleFunc("/upload/municipios", uploadHandler.UploadMunicipios).Methods("POST")
	api.HandleFunc("/upload/hospitais", uploadHandler.UploadHospitais).Methods("POST")
	api.HandleFunc("/upload/pacientes", uploadHandler.UploadPacientes).Methods("POST")
	api.HandleFunc("/upload/medicos", uploadHandler.UploadMedicos).Methods("POST")
	api.HandleFunc("/upload/cid10", uploadHandler.UploadCID10).Methods("POST")

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Printf("WebSocket endpoint: ws://localhost%s/ws", port)
	log.Fatal(http.ListenAndServe(port, router))
}