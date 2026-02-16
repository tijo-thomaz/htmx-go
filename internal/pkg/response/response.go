package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Responder handles HTTP responses
type Responder struct {
	log *slog.Logger
}

// New creates a new Responder
func New(log *slog.Logger) *Responder {
	return &Responder{
		log: log,
	}
}

// JSON sends a JSON response
func (r *Responder) JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		r.log.Error("json encoding failed", "error", err)
	}
}

// Error sends an error response
func (r *Responder) Error(w http.ResponseWriter, status int, message string) {
	r.log.Warn("error response", "status", status, "message", message)

	// For HTMX requests, return HTML error
	if w.Header().Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(status)
		w.Write([]byte(`<div class="error">` + message + `</div>`))
		return
	}

	// For API requests, return JSON
	r.JSON(w, status, map[string]string{"error": message})
}

// HXRedirect sends an HTMX redirect header
func (r *Responder) HXRedirect(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Redirect", url)
	w.WriteHeader(http.StatusOK)
}
