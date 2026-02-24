# Scene 8: Handler Setup & Auth (33:00 - 38:00)

> 🎬 **Previous**: Middleware + security (Scene 7)
> 🎯 **Goal**: Response helper, templates, handler wiring, auth (login/register/logout)

---

## Response Helper

**⌨️ Create `internal/pkg/response/response.go`:**
```go
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
	return &Responder{log: log}
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
	w.WriteHeader(status)
	w.Write([]byte(message))
}

// HXRedirect sends an HTMX redirect header
func (r *Responder) HXRedirect(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Redirect", url)
	w.WriteHeader(http.StatusOK)
}
```

> 🧠 📱 "HXRedirect — HTMX forms-ന് redirect ചെയ്യാൻ. Normal HTTP redirect HTMX ignore ചെയ്യും. HX-Redirect header set ചെയ്താൽ HTMX client-side redirect ചെയ്യും."
> 📱 "JSON() — API response-ന്. Content-Type set ചെയ്ത് JSON encode ചെയ്യും."

---

## Template Helper

**⌨️ Create `internal/pkg/templates/templates.go`:**
```go
package templates

import (
	"html/template"
	"io"
	"strings"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"multiply": func(a, b int) int { return a * b },
		"slice": func(s string, start, end int) string {
			if end > len(s) { end = len(s) }
			if start > len(s) { return "" }
			return s[start:end]
		},
		"upper": func(s string) string { return strings.ToUpper(s) },
	}
}

func Render(w io.Writer, page string, data interface{}) error {
	tmpl, err := template.New("base.html").Funcs(FuncMap()).ParseFiles(
		"web/templates/layouts/base.html",
		"web/templates/pages/"+page,
	)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(w, "base", data)
}
```

> 🧠 📱 "FuncMap — Go templates-ൽ custom functions. multiply AOS animation delay-ന്. slice avatar initial-ന്. upper capitalize-ന്."
> 📱 "Render — base layout + page template combine ചെയ്ത് execute."

---

## Handler Dependencies & Wiring

**⌨️ Create `internal/handler/handler.go`:**
```go
package handler

import (
	"log/slog"

	"linkbio/internal/pkg/response"
	"linkbio/internal/repository"

	"github.com/gorilla/sessions"
)

type Dependencies struct {
	Log           *slog.Logger
	Responder     *response.Responder
	Store         *sessions.CookieStore
	UserRepo      *repository.UserRepository
	LinkRepo      *repository.LinkRepository
	AnalyticsRepo *repository.AnalyticsRepository
}

type Handler struct {
	Health    *HealthHandler
	Auth      *AuthHandler
	Link      *LinkHandler
	Dashboard *DashboardHandler
	Profile   *ProfileHandler
}

func New(deps *Dependencies) *Handler {
	return &Handler{
		Health:    NewHealthHandler(deps),
		Auth:      NewAuthHandler(deps),
		Link:      NewLinkHandler(deps),
		Dashboard: NewDashboardHandler(deps),
		Profile:   NewProfileHandler(deps),
	}
}
```

> 🧠 📱 "Dependencies struct — dependency injection. എല്ലാ handlers-നും same dependencies. Testing-ന് easy mock."

---

## Health Handler

**⌨️ Create `internal/handler/health.go`:**
```go
package handler

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler(deps *Dependencies) *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
```

---

## Auth Handler — Login, Register, Logout

> 📱 "Authentication handler. bcrypt password hashing, session management, HTMX redirects."

**⌨️ Create `internal/handler/auth.go`:**
```go
package handler

import (
	"net/http"

	"linkbio/internal/model"
	"linkbio/internal/pkg/response"
	"linkbio/internal/pkg/templates"
	"linkbio/internal/repository"

	"log/slog"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	log      *slog.Logger
	resp     *response.Responder
	store    *sessions.CookieStore
	userRepo *repository.UserRepository
}

func NewAuthHandler(deps *Dependencies) *AuthHandler {
	return &AuthHandler{
		log:      deps.Log,
		resp:     deps.Responder,
		store:    deps.Store,
		userRepo: deps.UserRepo,
	}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if err := templates.Render(w, "login.html", nil); err != nil {
		h.log.Error("template error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.resp.Error(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		h.resp.Error(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	user, err := h.userRepo.GetByEmail(r.Context(), email)
	if err != nil {
		h.log.Error("database error", "error", err)
		h.resp.Error(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if user == nil {
		h.resp.Error(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		h.resp.Error(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Create session — gorilla auto-signs + encrypts the cookie!
	session, _ := h.store.Get(r, "session")
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	session.Save(r, w)

	h.log.Info("user logged in", "user_id", user.ID)
	h.resp.HXRedirect(w, "/dashboard")
}

func (h *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if err := templates.Render(w, "register.html", nil); err != nil {
		h.log.Error("template error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.resp.Error(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if username == "" || email == "" || password == "" {
		h.resp.Error(w, http.StatusBadRequest, "All fields are required")
		return
	}
	if len(password) < 6 {
		h.resp.Error(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	existing, _ := h.userRepo.GetByUsername(r.Context(), username)
	if existing != nil {
		h.resp.Error(w, http.StatusConflict, "Username already taken")
		return
	}
	existing, _ = h.userRepo.GetByEmail(r.Context(), email)
	if existing != nil {
		h.resp.Error(w, http.StatusConflict, "Email already registered")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		h.resp.Error(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		DisplayName:  username,
		Theme:        "light",
	}

	if err := h.userRepo.Create(r.Context(), user); err != nil {
		h.log.Error("user creation error", "error", err)
		h.resp.Error(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	session, _ := h.store.Get(r, "session")
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	session.Save(r, w)

	h.log.Info("user registered", "user_id", user.ID)
	h.resp.HXRedirect(w, "/dashboard")
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session")
	session.Options.MaxAge = -1 // Delete cookie
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
```

> 🧠 **Explain bcrypt:**
> 📱 "bcrypt — password hashing. Plain password database-ൽ store ചെയ്യരുത്! bcrypt.GenerateFromPassword hash create ചെയ്യും. bcrypt.CompareHashAndPassword verify ചെയ്യും."
> 📱 "Database leak ആയാലും original password recover ചെയ്യാൻ practically impossible."

> 🧠 **Explain session.Save:**
> 📱 "session.Save — gorilla internally cookie sign + encrypt (if key set) ചെയ്യും. നമ്മൾ plain values set ചെയ്താൽ മതി."

> 🧠 **Explain Logout:**
> 📱 "MaxAge = -1 — cookie delete ചെയ്യാൻ browser-നോട് പറയും. Session destroyed."

> 🎯 **Analogy:**
> 📱 "bcrypt = paper shredder. Password ഇട്ടാൽ hash വരും. Hash-ൽ നിന്ന് original password reconstruct ചെയ്യാൻ പറ്റില്ല. Same password ഇട്ടാൽ same hash — ഇങ്ങനെ verify ചെയ്യാം."

---

> 🎥 **Transition:** "Auth ready. ഇനി Link handler — CRUD + click tracking."
