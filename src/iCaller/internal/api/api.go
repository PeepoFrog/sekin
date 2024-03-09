package api

import (
	"encoding/json"
	"icaller/internal/command"
	"log"
	"net/http"
)

func ExecuteCommandHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Command string          `json:"command"`
		Args    json.RawMessage `json:"args"` // Use RawMessage for delayed parsing
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mapping, exists := command.CommandMapping[request.Command]
	if !exists {
		http.Error(w, "Command not allowed", http.StatusForbidden)
		return
	}

	// Decode the Args into the specific struct for this command
	args := mapping.ArgsStruct
	if err := json.Unmarshal(request.Args, &args); err != nil {
		http.Error(w, "Invalid arguments for command", http.StatusBadRequest)
		return
	}

	// Execute the command
	output, err := mapping.Handler(args)
	if err != nil {
		log.Printf("Error executing command '%s': %v", request.Command, err)
		http.Error(w, "Failed to execute command", http.StatusInternalServerError)
		return
	}
	log.Printf("Command '%s' executed successfully", request.Command)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"output": output})
}
