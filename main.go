package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	cfg := &apiConfig{
		fileserverHits: 0,
	}
	port := "8080"
	filepathRoot := "."
	sm := http.NewServeMux()

	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))
	sm.HandleFunc("GET /api/healthz", healthCheckHandler)
	sm.HandleFunc("GET /admin/metrics", cfg.getMetricsHandler)
	sm.HandleFunc("/api/reset", cfg.resetMetricsHandler)
	sm.HandleFunc("/api/validate_chirp", validateChirpHandler)
	sm.Handle("/app/*", cfg.middlewareMetricsInc(handler))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: sm,
	}

	log.Printf("Serving files from %s on port: %v\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
