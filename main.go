package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8080"
	filepathRoot := "."
	sm := http.NewServeMux()

	sm.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: sm,
	}

	log.Printf("Serving files from %s on port: %v\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
