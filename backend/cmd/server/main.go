package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Sebastian-Loza05/Sezzle-calculator/backend/internal/httpapi"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           httpapi.NewRouter(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("Starting server on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Could not listen on port %s: %v", port, err)
	}
}
