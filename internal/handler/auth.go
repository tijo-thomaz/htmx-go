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

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	log       *slog.Logger
	resp      *response.Responder
	store     *sessions.CookieStore
	userRepo  *repository.UserRepository
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(deps *Dependencies) *AuthHandler {
	return &AuthHandler{
		log:      deps.Log,
		resp:     deps.Responder,
		store:    deps.Store,
		userRepo: deps.UserRepo,
	}
}

// LoginPage renders the login page
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if err := templates.Render(w, "login.html", nil); err != nil {
		h.log.Error("template error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Login handles user login
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

	// Find user by email
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

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		h.resp.Error(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Create session
	session, _ := h.store.Get(r, "session")
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	if err := session.Save(r, w); err != nil {
		h.log.Error("session save error", "error", err)
		h.resp.Error(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	h.log.Info("user logged in", "user_id", user.ID, "username", user.Username)

	// HTMX redirect
	h.resp.HXRedirect(w, "/dashboard")
}

// RegisterPage renders the registration page
func (h *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if err := templates.Render(w, "register.html", nil); err != nil {
		h.log.Error("template error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Register handles user registration
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

	// Check if username exists
	existing, _ := h.userRepo.GetByUsername(r.Context(), username)
	if existing != nil {
		h.resp.Error(w, http.StatusConflict, "Username already taken")
		return
	}

	// Check if email exists
	existing, _ = h.userRepo.GetByEmail(r.Context(), email)
	if existing != nil {
		h.resp.Error(w, http.StatusConflict, "Email already registered")
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		h.log.Error("password hash error", "error", err)
		h.resp.Error(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Create user
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

	h.log.Info("user registered", "user_id", user.ID, "username", user.Username)

	// Create session
	session, _ := h.store.Get(r, "session")
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	session.Save(r, w)

	// HTMX redirect
	h.resp.HXRedirect(w, "/dashboard")
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
