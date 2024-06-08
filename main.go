package main

import (
	"log"
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func main() {
	port := "8080"
	filepathRoot := "."
	sm := http.NewServeMux()

	sm.HandleFunc("/healthz", healthCheckHandler)
	sm.Handle("/app/*", http.FileServer(http.Dir(filepathRoot)))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: sm,
	}

	log.Printf("Serving files from %s on port: %v\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
