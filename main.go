package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/WSMcReynolds/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	JWTSecret      string
}

func main() {
	godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")
	port := "8080"
	filepathRoot := "."
	path := "database.json"
	db, err := database.NewDB(path)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Started Database Connection at: %s", path)

	cfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		JWTSecret:      jwtSecret,
	}

	sm := http.NewServeMux()

	// File system Requests
	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))
	sm.Handle("/app/*", cfg.middlewareMetricsInc(handler))

	// Healthcheck endpoints
	sm.HandleFunc("GET /api/healthz", healthCheckHandler)

	// Admin endpionts
	sm.HandleFunc("GET /admin/metrics", cfg.metricsGetHandler)

	// Utility endpoints
	sm.HandleFunc("GET /api/reset", cfg.resetMetricsHandler)

	// Chirps endponits
	sm.HandleFunc("GET /api/chirps", cfg.chirpsRetrieveHandler)
	sm.HandleFunc("GET /api/chirps/{chirpID}", cfg.chirpsGetHandler)
	sm.HandleFunc("POST /api/chirps", cfg.chirpsCreateHandler)
	sm.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.chirpsDeleteHandler)

	// Users endpoints
	sm.HandleFunc("GET /api/users", cfg.usersRetrieveHandler)
	sm.HandleFunc("POST /api/users", cfg.usersCreateHandler)
	sm.HandleFunc("PUT /api/users", cfg.usersUpdateHandler)

	// Auth endpoints
	sm.HandleFunc("POST /api/login", cfg.loginRequestHandler)
	sm.HandleFunc("POST /api/refresh", cfg.refreshCreateHandler)
	sm.HandleFunc("POST /api/revoke", cfg.refreshRevokeHandler)

	// Webhook endpoints
	sm.HandleFunc("POST /api/polka/webhooks", cfg.webhookRequestHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: sm,
	}

	log.Printf("Serving files from %s on port: %v\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
