package main

import (
	"icaller/internal/api"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/execute", api.ExecuteCommandHandler)

	log.Println("Server is running on :8081...")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)

	}

}
