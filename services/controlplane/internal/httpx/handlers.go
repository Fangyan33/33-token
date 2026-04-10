package httpx

import (
	"encoding/json"
	"net/http"
	"time"
)

type placeholderResponse struct {
	Module  string `json:"module"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Healthz(w http.ResponseWriter, _ *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func Placeholder(module string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		WriteJSON(w, http.StatusOK, placeholderResponse{
			Module:  module,
			Status:  "bootstrap_ready",
			Message: "module scaffold created, implementation pending",
		})
	}
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
