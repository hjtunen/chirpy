package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
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

	respondWithJSON(w, http.StatusOK, map[string]bool{"valid": true})

}

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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	errorResponse := map[string]string{"error": message}
	jsonData, err := json.Marshal(errorResponse)
	if err != nil {
		log.Printf("Error marshaling JSON error response: %v", err)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(jsonData)
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	appHandler := apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", http.StripPrefix("/app/", appHandler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Starting server on port %s, serving files from %s", port, filepathRoot)
	log.Fatal(server.ListenAndServe())

}
