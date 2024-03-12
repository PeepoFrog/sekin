package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sekinEngine/internal/api"
)

func main() {
	logFilePath := filepath.Join("/sekinEngine", "sekinEngine.log")

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	const port int = 9001
	// Log that the server is starting.
	log.Printf("Server is starting on :%v...", port)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/execute", api.ExecuteCommandHandler)

	log.Printf("Server is running on :%v...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)

	}

}
