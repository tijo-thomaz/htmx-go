# Scene 11: Router & Server (45:00 - 48:00)

> 🎬 **Previous**: All handlers done (Scene 10)
> 🎯 **Goal**: Wire routes (public vs protected), server setup
> ⚠️ **Teaching highlight**: Click route placement — public, not behind auth!

---

## Router — ⚠️ Public vs Protected Routes

> 📱 "Router-ൽ routes define ചെയ്യാം. ⚠️ ഏത് routes public, ഏത് protected — ഈ decision critical."

**⌨️ Create `internal/router/router.go`:**
```go
package router

import (
	"net/http"

	"linkbio/internal/handler"
	"linkbio/internal/middleware"
	"linkbio/internal/pkg/templates"

	"github.com/go-chi/chi/v5"
)

func New(h *handler.Handler, mw *middleware.Middleware) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware — every request goes through these
	r.Use(mw.Recovery)   // Catch panics
	r.Use(mw.Logger)     // Log requests
	r.Use(mw.RateLimit)  // Prevent abuse

	r.Get("/health", h.Health.Check)

	fileServer := http.FileServer(http.Dir("web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// ✅ PUBLIC routes — no login required
	r.Group(func(r chi.Router) {
		r.Get("/", handleHome)
		r.Get("/u/{username}", h.Profile.Show)
		r.Get("/click/{id}", h.Link.Click) // ⚠️ PUBLIC! Not behind auth!
	})

	// Auth pages
	r.Route("/auth", func(r chi.Router) {
		r.Get("/login", h.Auth.LoginPage)
		r.Post("/login", h.Auth.Login)
		r.Get("/register", h.Auth.RegisterPage)
		r.Post("/register", h.Auth.Register)
		r.Post("/logout", h.Auth.Logout)
	})

	// 🔒 PROTECTED API — requires login
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(mw.Auth)
		r.Route("/links", func(r chi.Router) {
			r.Post("/", h.Link.Create)
			r.Put("/{id}", h.Link.Update)
			r.Delete("/{id}", h.Link.Delete)
			r.Post("/reorder", h.Link.Reorder)
		})
	})

	// 🔒 PROTECTED Dashboard
	r.Route("/dashboard", func(r chi.Router) {
		r.Use(mw.Auth)
		r.Get("/", h.Dashboard.Index)
	})

	return r
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if err := templates.Render(w, "home.html", nil); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
```

> 🧠 **⚠️ Why `/click/{id}` is in the PUBLIC group:**
> 📱 "ഈ line ആണ് click tracking work ചെയ്യുന്നതിന്റെ key."
> 📱 "Profile page visit ചെയ്യുന്ന visitors logged in അല്ല. Link click ചെയ്യുമ്പോൾ GET /click/5 request."
> 📱 "ഇത് auth group-ൽ ഇട്ടാൽ → middleware redirect to /auth/login → click lost → dashboard 0 clicks."
> 📱 "GET method ആണ്, POST അല്ല. `<a href='/click/5'>` standard browser GET."

> 🧠 **Explain middleware groups:**
> 📱 "chi.Router Group — same middleware share ചെയ്യുന്ന routes."
> 📱 "r.Use(mw.Auth) — ആ group-ലെ എല്ലാ routes-നും auth check."
> 📱 "Public group-ൽ Auth middleware ഇല്ല — anyone access ചെയ്യാം."

---

## Server Setup

**⌨️ Create `internal/server/server.go`:**
```go
package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"linkbio/internal/config"
	"linkbio/internal/handler"
	"linkbio/internal/middleware"
	"linkbio/internal/pkg/response"
	"linkbio/internal/repository"
	"linkbio/internal/router"
)

type Server struct {
	httpServer *http.Server
	log        *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) (*Server, error) {
	db, err := repository.NewDB(cfg.DatabasePath, log)
	if err != nil {
		return nil, err
	}
	if err := repository.Migrate(db, log); err != nil {
		return nil, err
	}

	userRepo := repository.NewUserRepository(db)
	linkRepo := repository.NewLinkRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)

	resp := response.New(log)

	// ⚠️ 4 params: log, signingKey, encryptionKey, rateLimit
	mw := middleware.New(log, cfg.SessionSecret, cfg.SessionEncKey, cfg.RateLimit)

	h := handler.New(&handler.Dependencies{
		Log:           log,
		Responder:     resp,
		Store:         mw.Store(),
		UserRepo:      userRepo,
		LinkRepo:      linkRepo,
		AnalyticsRepo: analyticsRepo,
	})

	r := router.New(h, mw)

	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{httpServer: httpServer, log: log}, nil
}

func (s *Server) Start() error {
	s.log.Info("server starting", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("server shutting down")
	return s.httpServer.Shutdown(ctx)
}
```

> 🧠 📱 "middleware.New — 4 parameters: log, signingKey, encryptionKey, rateLimit."
> 📱 "cfg.SessionEncKey — empty ആയാൽ signing only + warning. 32 bytes ആയാൽ AES-256 encryption."
> 📱 "ReadTimeout, WriteTimeout, IdleTimeout — production security. Slow clients server block ചെയ്യില്ല."

---

> 🎥 **Transition:** "Backend complete! ഇനി frontend — templates."
