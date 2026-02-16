package middleware

import (
	"net/http"
)

// RateLimit limits the number of requests per second
func (m *Middleware) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.limiter.Allow() {
			m.log.Warn("rate limit exceeded",
				"ip", r.RemoteAddr,
				"path", r.URL.Path,
			)
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
