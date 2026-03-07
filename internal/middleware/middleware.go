package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
)

// Middleware holds dependencies for all middleware functions
type Middleware struct {
	log     *slog.Logger
	store   *sessions.CookieStore
}

// New creates a new Middleware instance
func New(log *slog.Logger, sessionSecret, encryptionKey string) *Middleware {
	var store *sessions.CookieStore
	if len(encryptionKey) == 32 {
		store = sessions.NewCookieStore([]byte(sessionSecret), []byte(encryptionKey))
		log.Info("session cookies: signing + AES-256 encryption enabled")
	} else {
		store = sessions.NewCookieStore([]byte(sessionSecret))
		log.Warn("session cookies: signing only (set SESSION_ENCRYPTION_KEY for encryption)")
	}
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	return &Middleware{
		log:   log,
		store: store,
	}
}

// Store returns the session store
func (m *Middleware) Store() *sessions.CookieStore {
	return m.store
}
