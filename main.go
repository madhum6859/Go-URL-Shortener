package main

import (
	"log"
	"net/http"

	"c:\Users\iamma\Programs\trae\Go-URL-Shortener\handlers"
	"c:\Users\iamma\Programs\trae\Go-URL-Shortener\storage"
)

func main() {
	// Initialize storage
	store := storage.NewInMemoryStore()

	// Initialize handlers
	handler := handlers.NewHandler(store)

	// Set up routes
	http.HandleFunc("/", handler.RedirectHandler)
	http.HandleFunc("/shorten", handler.ShortenHandler)
	http.HandleFunc("/health", handler.HealthCheckHandler)

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}