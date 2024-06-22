package main

import (
	"log"
	"net/http"

	"github.com/WSMcReynolds/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
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
	}

	sm := http.NewServeMux()

	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))
	sm.HandleFunc("GET /api/healthz", healthCheckHandler)
	sm.HandleFunc("GET /admin/metrics", cfg.metricsGetHandler)
	sm.HandleFunc("GET /api/reset", cfg.resetMetricsHandler)
	sm.HandleFunc("GET /api/chirps", cfg.chirpsRetrieveHandler)
	sm.HandleFunc("POST /api/chirps", cfg.chirpsCreateHandler)
	sm.HandleFunc("GET /api/chirps/{chirpID}", cfg.chirpsGetHandler)
	sm.HandleFunc("GET /api/users", cfg.usersRetrieveHandler)
	sm.HandleFunc("POST /api/users", cfg.usersCreateHandler)
	sm.HandleFunc("POST /api/login", cfg.loginRequestHandler)
	sm.Handle("/app/*", cfg.middlewareMetricsInc(handler))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: sm,
	}

	log.Printf("Serving files from %s on port: %v\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
