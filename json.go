package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON response: %v", err)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(jsonData)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string, err error) {
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Error: %s", message)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, statusCode, errorResponse{Error: message})
}
