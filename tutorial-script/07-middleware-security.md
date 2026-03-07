# Scene 7: Middleware & Security Deep Dive (28:00 - 33:00)

> 🎬 **Previous**: Database layer complete (Scene 6)
> 🎯 **Goal**: Build middleware chain + understand session encryption + why not JWT

---

## 🎥 Transition

> 📱 "Database ready. ഇനി middleware — request handler-ന് എത്തുന്നതിന് മുമ്പ് run ആകുന്ന code."

> 🎯 **Analogy:**
> 📱 "Airport പോലെ. Plane-ൽ കയറും മുമ്പ് security check, ticket check, boarding pass check. ഓരോന്നും ഒരു middleware."

---

## Core Middleware — Session Store with Encryption

**⌨️ Create `internal/middleware/middleware.go`:**
```go
package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
)

type Middleware struct {
	log     *slog.Logger
	store   *sessions.CookieStore
}

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
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	return &Middleware{
		log:     log,
		store:   store,
	}
}

func (m *Middleware) Store() *sessions.CookieStore {
	return m.store
}
```

> 🧠 **Explain — ⚠️ This is the KEY security feature:**

> 📱 "sessions.NewCookieStore-ന് രണ്ട് keys pass ചെയ്യാം."
> 📱 "ഒരു key = signing only. Tamper ചെയ്യാൻ പറ്റില്ല, but Base64 decode ചെയ്ത് user_id read ചെയ്യാം!"
> 📱 "രണ്ട് keys = signing + AES-256 encryption. Read-ഉം modify-ഉം ചെയ്യാൻ പറ്റില്ല!"

> 🎯 **Seal vs Locked Box analogy:**
> 📱 "Signing = letter-ന് wax seal. Break ചെയ്യാതെ change ചെയ്യാൻ പറ്റില്ല. But letter read ചെയ്യാം — transparent envelope."
> 📱 "Encryption = letter locked steel box-ൽ. Key ഇല്ലാതെ open ചെയ്യാനോ read ചെയ്യാനോ പറ്റില്ല."
> 📱 "നമുക്ക് രണ്ടും വേണം — sealed AND locked!"

> 📱 "len(encryptionKey) == 32 check — AES-256-ന് exactly 32 bytes വേണം. Wrong size ആയാൽ encryption skip, warning log."
> 📱 "gorilla/securecookie internally AES-CTR + HMAC-SHA256 handle ചെയ്യും. നമ്മൾ key pass ചെയ്താൽ മതി."

> 📱 "Cookie options:"
> 📱 "HttpOnly = true — JavaScript-ന് cookie access ഇല്ല. XSS protection."
> 📱 "SameSite Lax — CSRF protection. Third-party sites-ൽ നിന്ന് cookie send ചെയ്യില്ല."
> 📱 "Secure = false — development. Production-ൽ true set ചെയ്യണം (HTTPS only)."

---

## Logger Middleware

**⌨️ Create `internal/middleware/logger.go`:**
```go
package middleware

import (
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (m *Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		m.log.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.status,
			"duration", time.Since(start).String(),
			"ip", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)
	})
}
```

> 🧠 📱 "ResponseWriter wrap ചെയ്യുന്നു status code capture ചെയ്യാൻ. Default Go ResponseWriter status code expose ചെയ്യില്ല."

---

## Recovery Middleware

**⌨️ Create `internal/middleware/recovery.go`:**
```go
package middleware

import (
	"net/http"
	"runtime/debug"
)

func (m *Middleware) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.log.Error("panic recovered",
					"error", err,
					"path", r.URL.Path,
					"stack", string(debug.Stack()),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
```

> 🧠 📱 "Panic = Go crash. defer recover() catch ചെയ്ത് 500 error response. Server crash ആകില്ല."

---

## Auth Middleware

**⌨️ Create `internal/middleware/auth.go`:**
```go
package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UsernameKey contextKey = "username"
)

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
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/auth/login")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		username, _ := session.Values["username"].(string)

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UsernameKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) int64 {
	if id, ok := ctx.Value(UserIDKey).(int64); ok {
		return id
	}
	return 0
}

func UsernameFromContext(ctx context.Context) string {
	if name, ok := ctx.Value(UsernameKey).(string); ok {
		return name
	}
	return ""
}
```

> 🧠 📱 "Session cookie read → user_id extract → context-ൽ inject → handler-ൽ available."
> 📱 "HX-Request header check — HTMX request ആണെങ്കിൽ HX-Redirect header. Browser request ആണെങ്കിൽ HTTP 303 redirect."
> 📱 "context.WithValue — request-ന്റെ കൂടെ data pass. Middleware set, handler read."

---

## 🔐 Security Deep Dive: Why Sessions, Not JWT?

> 🎥 **Camera note:** Pause coding. Talk directly to camera. This is a teaching moment.

> 📱 "ഇവിടെ ഒരു important topic. Authentication-ന് JWT use ചെയ്യണോ Sessions use ചെയ്യണോ? YouTube tutorials-ൽ JWT everywhere കാണും. But production-ൽ..."

### ⚠️ Why Industries Don't Trust JWT for User Sessions

| Problem | Explanation |
|---------|-------------|
| **Revoke ചെയ്യാൻ പറ്റില്ല** | Issue ചെയ്ത JWT expire ആകുന്നത് വരെ valid. User hack ആയാൽ? Blocklist build ചെയ്യണം — "stateless" advantage gone! |
| **Token theft = full access** | localStorage-ൽ store → XSS steal. Cookie-ൽ store → stateless advantage gone. |
| **Size bloat** | Session cookie = ~32 bytes. JWT = ~800+ bytes. ഓരോ request-ലും extra bandwidth. |
| **Crypto footguns** | `alg: none` attack, RS256/HS256 confusion. Spec too flexible, bugs too easy. |
| **"Stateless" is a myth** | Practice-ൽ: user deleted? permissions changed? password reset? Database check avoid ചെയ്യാൻ പറ്റില്ല. |

### 🎯 The Correct Approach

| Use Case | Correct Solution |
|----------|-----------------|
| **User sessions** | Signed + encrypted cookies (GitHub, Stripe, banks) |
| **Service-to-service API** | Short-lived JWTs (5-15 min), OAuth2 server |
| **Third-party login** | OAuth2/OIDC (Google, GitHub login) |

> 📱 "നമ്മുടെ app HTMX ആണ്. ഓരോ request server hit ചെയ്യും. Session cookies = correct pattern. Simple, secure, revocable."

> 🎯 📱 "JWT bad അല്ല — wrong use case-ന് use ചെയ്യുന്നത് ആണ് problem. Screwdriver excellent tool ആണ്. Nail അടിക്കാൻ use ചെയ്യരുത്."

> 🎥 **Transition:** "Security clear. ഇനി handlers — actual business logic!"
