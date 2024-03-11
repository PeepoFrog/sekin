package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"scaller/internal/api"
)

func main() {
	logFilePath := filepath.Join("/sekai", "scaller.log")

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	// Log that the server is starting.
	log.Println("Server is starting on :8080...")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/execute", api.ExecuteCommandHandler)

	log.Println("Server is running on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)

	}

}
