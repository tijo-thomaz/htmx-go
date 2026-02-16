package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"linkbio/internal/model"
	"linkbio/internal/pkg/response"
	"linkbio/internal/repository"
	"linkbio/internal/testutil"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func setupAuthHandler(t *testing.T) (*AuthHandler, *repository.UserRepository) {
	t.Helper()

	db := testutil.TestDB(t)
	log := testutil.TestLogger()

	userRepo := repository.NewUserRepository(db)
	resp := response.New(log)
	store := sessions.NewCookieStore([]byte("test-secret-key-32-chars-minimum!"))

	h := &AuthHandler{
		log:      log,
		resp:     resp,
		store:    store,
		userRepo: userRepo,
	}

	return h, userRepo
}

func TestAuthHandler_Login_Success(t *testing.T) {
	h, userRepo := setupAuthHandler(t)
	ctx := context.Background()

	// Setup: Create user with hashed password
	password := "testpass123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &model.User{
		Username:     "loginuser",
		Email:        "login@test.com",
		PasswordHash: string(hash),
		DisplayName:  "Login User",
		Theme:        "light",
	}
	userRepo.Create(ctx, user)

	// Create login request
	form := url.Values{}
	form.Add("email", "login@test.com")
	form.Add("password", password)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	// Check: Should redirect via HX-Redirect header
	if rec.Header().Get("HX-Redirect") != "/dashboard" {
		t.Errorf("HX-Redirect = %s, want /dashboard", rec.Header().Get("HX-Redirect"))
	}
}

func TestAuthHandler_Login_WrongPassword(t *testing.T) {
	h, userRepo := setupAuthHandler(t)
	ctx := context.Background()

	// Setup: Create user
	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
	user := &model.User{
		Username:     "wrongpass",
		Email:        "wrong@test.com",
		PasswordHash: string(hash),
		DisplayName:  "Wrong Pass",
		Theme:        "light",
	}
	userRepo.Create(ctx, user)

	// Try login with wrong password
	form := url.Values{}
	form.Add("email", "wrong@test.com")
	form.Add("password", "wrongpassword")

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	// Check: Should return 401 Unauthorized
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthHandler_Login_UserNotFound(t *testing.T) {
	h, _ := setupAuthHandler(t)

	form := url.Values{}
	form.Add("email", "notexists@test.com")
	form.Add("password", "anypassword")

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthHandler_Login_MissingFields(t *testing.T) {
	h, _ := setupAuthHandler(t)

	tests := []struct {
		name     string
		email    string
		password string
	}{
		{"missing email", "", "password"},
		{"missing password", "test@test.com", ""},
		{"missing both", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("email", tt.email)
			form.Add("password", tt.password)

			req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()

			h.Login(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("Status = %d, want %d", rec.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestAuthHandler_Register_Success(t *testing.T) {
	h, userRepo := setupAuthHandler(t)
	ctx := context.Background()

	form := url.Values{}
	form.Add("username", "newuser")
	form.Add("email", "new@test.com")
	form.Add("password", "password123")

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Register(rec, req)

	// Check redirect
	if rec.Header().Get("HX-Redirect") != "/dashboard" {
		t.Errorf("HX-Redirect = %s, want /dashboard", rec.Header().Get("HX-Redirect"))
	}

	// Verify user was created
	user, _ := userRepo.GetByUsername(ctx, "newuser")
	if user == nil {
		t.Error("User was not created")
	}

	// Verify password was hashed
	if user.PasswordHash == "password123" {
		t.Error("Password was not hashed")
	}
}

func TestAuthHandler_Register_DuplicateUsername(t *testing.T) {
	h, userRepo := setupAuthHandler(t)
	ctx := context.Background()

	// Create existing user
	existing := &model.User{
		Username:     "existing",
		Email:        "existing@test.com",
		PasswordHash: "hash",
		DisplayName:  "Existing",
		Theme:        "light",
	}
	userRepo.Create(ctx, existing)

	// Try to register with same username
	form := url.Values{}
	form.Add("username", "existing")
	form.Add("email", "different@test.com")
	form.Add("password", "password123")

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Register(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusConflict)
	}
}

func TestAuthHandler_Register_DuplicateEmail(t *testing.T) {
	h, userRepo := setupAuthHandler(t)
	ctx := context.Background()

	// Create existing user
	existing := &model.User{
		Username:     "existing",
		Email:        "existing@test.com",
		PasswordHash: "hash",
		DisplayName:  "Existing",
		Theme:        "light",
	}
	userRepo.Create(ctx, existing)

	// Try to register with same email
	form := url.Values{}
	form.Add("username", "newuser")
	form.Add("email", "existing@test.com")
	form.Add("password", "password123")

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Register(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusConflict)
	}
}

func TestAuthHandler_Register_ShortPassword(t *testing.T) {
	h, _ := setupAuthHandler(t)

	form := url.Values{}
	form.Add("username", "shortpass")
	form.Add("email", "short@test.com")
	form.Add("password", "12345") // Less than 6 chars

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestAuthHandler_Register_MissingFields(t *testing.T) {
	h, _ := setupAuthHandler(t)

	tests := []struct {
		name     string
		username string
		email    string
		password string
	}{
		{"missing username", "", "test@test.com", "password"},
		{"missing email", "testuser", "", "password"},
		{"missing password", "testuser", "test@test.com", ""},
		{"missing all", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("username", tt.username)
			form.Add("email", tt.email)
			form.Add("password", tt.password)

			req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()

			h.Register(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("Status = %d, want %d", rec.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	h, _ := setupAuthHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	rec := httptest.NewRecorder()

	h.Logout(rec, req)

	// Check: Should redirect to home
	if rec.Code != http.StatusSeeOther {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusSeeOther)
	}

	location := rec.Header().Get("Location")
	if location != "/" {
		t.Errorf("Location = %s, want /", location)
	}
}

// Benchmarks

// Note: Benchmark for auth is slow due to bcrypt, skipped in normal runs
// Run with: go test -bench=. -run=^$ ./internal/handler
