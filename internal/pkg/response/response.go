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

// Error sends an error response (always HTML for simplicity since all clients are HTMX/browser)
func (r *Responder) Error(w http.ResponseWriter, status int, message string) {
	r.log.Warn("error response", "status", status, "message", message)
	w.WriteHeader(status)
	w.Write([]byte(message))
}

// HXRedirect sends an HTMX redirect header
func (r *Responder) HXRedirect(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Redirect", url)
	w.WriteHeader(http.StatusOK)
}
