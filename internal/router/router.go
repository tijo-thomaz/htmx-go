package router

import (
	"net/http"

	"linkbio/internal/handler"
	"linkbio/internal/middleware"
	"linkbio/internal/pkg/templates"

	"github.com/go-chi/chi/v5"
)

// New creates the main router with all routes
func New(h *handler.Handler, mw *middleware.Middleware) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware chain
	r.Use(mw.Recovery)   // Recover from panics
	r.Use(mw.Logger)     // Log all requests
	r.Use(mw.RateLimit)  // Rate limiting

	// Health check (no auth required)
	r.Get("/health", h.Health.Check)

	// Static files
	fileServer := http.FileServer(http.Dir("web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", handleHome)
		r.Get("/u/{username}", h.Profile.Show)
		r.Get("/click/{id}", h.Link.Click)
	})

	// Auth namespace
	r.Route("/auth", func(r chi.Router) {
		r.Get("/login", h.Auth.LoginPage)
		r.Post("/login", h.Auth.Login)
		r.Get("/register", h.Auth.RegisterPage)
		r.Post("/register", h.Auth.Register)
		r.Post("/logout", h.Auth.Logout)
	})

	// API v1 namespace (protected)
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(mw.Auth) // Require authentication

		r.Route("/links", func(r chi.Router) {
			r.Post("/", h.Link.Create)
			r.Put("/{id}", h.Link.Update)
			r.Delete("/{id}", h.Link.Delete)
			r.Post("/reorder", h.Link.Reorder)
		})


	})

	// Dashboard namespace (protected)
	r.Route("/dashboard", func(r chi.Router) {
		r.Use(mw.Auth)
		r.Get("/", h.Dashboard.Index)
	})

	return r
}

// handleHome renders the landing page
func handleHome(w http.ResponseWriter, r *http.Request) {
	if err := templates.Render(w, "home.html", nil); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
