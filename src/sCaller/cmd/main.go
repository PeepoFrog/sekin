package main

import (
	"log"
	"net/http"
	"scaller/internal/api"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/execute", api.ExecuteCommandHandler)

	log.Println("Server is running on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)

	}

}
