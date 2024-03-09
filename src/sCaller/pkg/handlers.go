package handlers

import (
	"encoding/json"
	"net/http"
)

func ResourceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Handel GET requests
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode("GET request to /api/resource")
	case "POST":
		// Handle POST requests
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode("POST request to /api/resource")
	default:
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode("POST request to /api/resource")
	}
}
