package middleware

import (
	"context"
	"net/http"
)

// Context keys
type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UsernameKey contextKey = "username"
)

// Auth checks if user is authenticated
func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := m.store.Get(r, "session")
		if err != nil {
			m.log.Error("session error", "error", err)
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		userID, ok := session.Values["user_id"].(int64)
		if !ok || userID == 0 {
			// For HTMX requests, return redirect header
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/auth/login")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		username, _ := session.Values["username"].(string)

		// Add user info to context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UsernameKey, username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserIDFromContext retrieves user ID from context
func UserIDFromContext(ctx context.Context) int64 {
	userID, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		return 0
	}
	return userID
}

// UsernameFromContext retrieves username from context
func UsernameFromContext(ctx context.Context) string {
	username, ok := ctx.Value(UsernameKey).(string)
	if !ok {
		return ""
	}
	return username
}
