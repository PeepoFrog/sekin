package main

import (
	"icaller/internal/api"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	logFilePath := filepath.Join("/interx", "icaller.log")

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	log.Println("Server is starting on :8081...")
	log.SetOutput(logFile)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/execute", api.ExecuteCommandHandler)

	log.Println("Server is running on :8081...")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)

	}

}
