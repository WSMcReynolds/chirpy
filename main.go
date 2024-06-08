package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8080"
	sm := http.NewServeMux()

	server := http.Server{
		Addr:    ":" + port,
		Handler: sm,
	}

	sm.Handle("/", http.FileServer(http.Dir(".")))

	log.Printf("Serving on port: %v\n", port)
	log.Fatal(server.ListenAndServe())
}
