package handler

import (
	"log/slog"
	"net/http"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	log *slog.Logger
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(log *slog.Logger) *HealthHandler {
	return &HealthHandler{log: log}
}

// Check returns the health status
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
