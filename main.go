package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "pawit-booking-bff"})
	})
	mux.HandleFunc("GET /api/v1/public/slots", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"items": []map[string]string{
			{"id": "slot_001", "clinic": "PawIt VetCare", "vet": "Dr. Asha Rao", "time": "09:30", "type": "Check-up"},
		}})
	})
	slog.Info("PawIt Booking BFF started", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
