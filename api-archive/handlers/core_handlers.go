package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ServerHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type Response struct {
		Message string `json:"message"`
	}

	var response Response

	if serverIsHealthy() {
		w.WriteHeader(http.StatusOK)
		response.Message = "Server is healthy"
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = "Server is not healthy"
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

func serverIsHealthy() bool {
	// Check the health of the server and return true or false accordingly
	// For example, check if the server can connect to the database
	return true
}
