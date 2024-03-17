package commands

import (
	"encoding/json"
	"net/http"
)

type CommandRequest struct {
	Command string                 `json:"command"`
	Args    map[string]interface{} `json:"args"`
}

func ExecuteCommandHandler(w http.ResponseWriter, r *http.Request) {
	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	handler, exists := GetCommandHandler(req.Command)
	if !exists {
		http.Error(w, "Command not supported", http.StatusNotFound)
		return
	}

	if err := handler.HandleCommand(req.Args); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
