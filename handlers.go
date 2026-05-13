package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hjtunen/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpsOne(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("chirpID")
	chirp, err := cfg.db.GetChirp(req.Context(), uuid.MustParse(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Error fetching chirp: %v", err)
		return
	}

	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	jsonChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, jsonChirp)
}

func (cfg *apiConfig) handlerChirpsList(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.ListChirps(req.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error listing chirps: %v", err)
		return
	}

	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	jsonChirps := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		jsonChirps[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, jsonChirps)
}

func (cfg *apiConfig) handlerChirpCreate(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
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

	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		UserID: param.UserID,
		Body:   cleanedBody,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error creating chirp: %v", err)
		return
	}

	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	jsonChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, jsonChirp)
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	param := parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	user, err := cfg.db.CreateUser(req.Context(), param.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error creating user: %v", err)
		return
	}

	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	jsonUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusCreated, jsonUser)
}

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
