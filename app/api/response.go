package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func OKResponse(w http.ResponseWriter, data any) {
	OKResponseWithStatus(w, data, http.StatusOK)
}

func OKResponseWithStatus(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("ERROR: Failed to encode JSON response: %v", err)

	}
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResponse := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Printf("ERROR: Failed to encode JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

}
