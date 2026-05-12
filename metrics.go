package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	payload := fmt.Sprintf(
		`<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(payload))
}
func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, req *http.Request) {

	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Resetting metrics is only allowed in dev environment", nil)
		return
	}

	cfg.fileserverHits.Store(0)

	error := cfg.db.DeleteAllUsers(req.Context())
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error deleting users: %v", error)
		return
	}

	type returnVals struct {
		Message string `json:"message"`
	}

	respondWithJSON(w, http.StatusOK, returnVals{Message: "Metrics reset successfully"})

}
