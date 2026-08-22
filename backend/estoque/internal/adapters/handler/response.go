package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type errorResponse struct {
	Status    int      `json:"status"`
	Message   string   `json:"message"`
	Timestamp string   `json:"timestamp"`
	Details   []string `json:"details"`
}

func decodeJSON(r *http.Request, destination any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(destination)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string, details ...string) {
	writeJSON(w, status, errorResponse{Status: status, Message: message, Timestamp: time.Now().UTC().Format(time.RFC3339), Details: details})
}

func logUnexpected(context string, err error) {
	if err != nil {
		log.Printf("%s: %v", context, err)
	}
}
