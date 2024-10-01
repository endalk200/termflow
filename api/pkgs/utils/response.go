package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func Response(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
	}
}

func ResponseError(w http.ResponseWriter, statusCode int, msg string) {
	if statusCode > 499 {
		log.Println("Reponding with 5XX error", msg)
	}

	type errResponse struct {
		Error string `json:"error"`
	}

	Response(w, statusCode, errResponse{
		Error: msg,
	})
}
