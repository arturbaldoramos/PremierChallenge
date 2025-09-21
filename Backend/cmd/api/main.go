package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"example.com/m/v2/internal/config"
	"example.com/m/v2/internal/database"
	"example.com/m/v2/internal/handlers"
	"example.com/m/v2/internal/services"
	"github.com/gorilla/mux"
)

// getEnvAsInt retrieves an environment variable as an integer, with a default fallback
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

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
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	redisAddr := redisHost + ":" + redisPort
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := getEnvAsInt("REDIS_DB", 0)

	log.Printf("Connecting to Redis at: %s", redisAddr)
	redisService := services.NewRedisService(redisAddr, redisPassword, redisDB)
	dataService := services.NewDataService(db)
	workerService := services.NewWorkerService(redisService, dataService)
	monitorService := services.NewMonitorService(redisService)

	// Initialize handlers
	uploadHandler := handlers.NewUploadHandler(db)
	wsHandler := handlers.NewWebSocketHandler(redisService, dataService)
	statsHandler := handlers.NewStatsHandler(db)

	// Start background services
	workerConfig := services.WorkerConfig{
		MaxWorkers: map[string]int{
			"estados":    getEnvAsInt("MAX_WORKERS_ESTADOS", 2),
			"municipios": getEnvAsInt("MAX_WORKERS_MUNICIPIOS", 4),
			"medicos":    getEnvAsInt("MAX_WORKERS_MEDICOS", 3),
			"hospitais":  getEnvAsInt("MAX_WORKERS_HOSPITAIS", 2),
			"pacientes":  getEnvAsInt("MAX_WORKERS_PACIENTES", 3),
			"cid10":      getEnvAsInt("MAX_WORKERS_CID10", 2),
		},
		BatchSize: map[string]int{
			"estados":    getEnvAsInt("BATCH_SIZE_ESTADOS", 50),
			"municipios": getEnvAsInt("BATCH_SIZE_MUNICIPIOS", 200),
			"medicos":    getEnvAsInt("BATCH_SIZE_MEDICOS", 300),
			"hospitais":  getEnvAsInt("BATCH_SIZE_HOSPITAIS", 150),
			"pacientes":  getEnvAsInt("BATCH_SIZE_PACIENTES", 250),
			"cid10":      getEnvAsInt("BATCH_SIZE_CID10", 400),
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

	// Statistics and filtering routes
	stats := api.PathPrefix("/stats").Subrouter()
	stats.HandleFunc("/totals", statsHandler.GetTotalStats).Methods("GET")
	stats.HandleFunc("/hospitais-por-especialidade", statsHandler.GetHospitalsBySpecialty).Methods("GET")
	stats.HandleFunc("/hospitais-por-municipio", statsHandler.GetHospitalsByMunicipio).Methods("GET")
	stats.HandleFunc("/medicos-distribuicao", statsHandler.GetMedicosDistribution).Methods("GET")
	stats.HandleFunc("/especialidades", statsHandler.GetEspecialidades).Methods("GET")
	stats.HandleFunc("/cid10-mais-comuns", statsHandler.GetCid10MaisComuns).Methods("GET")
	stats.HandleFunc("/hospitais-mais-acessados", statsHandler.GetHospitaisMaisAcessados).Methods("GET")
	stats.HandleFunc("/debug-cid10", statsHandler.GetDebugCid10).Methods("GET")
	stats.HandleFunc("/assign-medicos-hospitais", statsHandler.AssignMedicosToHospitals).Methods("POST")

	// Stats2 routes - Nova rota para estados e municípios
	api.HandleFunc("/stats2", statsHandler.GetStats2).Methods("GET")

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Printf("WebSocket endpoint: ws://localhost%s/ws", port)
	log.Fatal(http.ListenAndServe(port, router))
}
