package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("OK"))
}

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	param := parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	if len(param.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedBody := removeBadWords(param.Body)

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	respondWithJSON(w, http.StatusOK, returnVals{CleanedBody: cleanedBody})
}

func removeBadWords(input string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(input, " ")

	for i, word := range words {
		for _, badWord := range badWords {
			if strings.EqualFold(word, badWord) {
				words[i] = "****"
				break
			}
		}
	}

	cleaned := strings.Join(words, " ")

	return cleaned
}
