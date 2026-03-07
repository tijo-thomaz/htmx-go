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
		r.Get("/stats", h.Dashboard.Stats) // HTMX polling endpoint
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

	// ⚠️ 3 params: log, signingKey, encryptionKey
	mw := middleware.New(log, cfg.SessionSecret, cfg.SessionEncKey)

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

> 🧠 📱 "middleware.New — 3 parameters: log, signingKey, encryptionKey."
> 📱 "cfg.SessionEncKey — empty ആയാൽ signing only + warning. 32 bytes ആയാൽ AES-256 encryption."
> 📱 "ReadTimeout, WriteTimeout, IdleTimeout — production security. Slow clients server block ചെയ്യില്ല."

---

## Entry Point — main.go

> 📱 "Server struct ready. ഇനി actual entry point — main.go. Application start, config load, graceful shutdown — എല്ലാം ഇവിടെ."

**⌨️ Create `cmd/server/main.go`:**
```go
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"linkbio/internal/config"
	"linkbio/internal/pkg/logger"
	"linkbio/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Initialize logger
	var log = logger.New(cfg.LogLevel)
	if cfg.IsDevelopment() {
		log = logger.NewDevelopment()
	}

	log.Info("starting linkbio",
		"env", cfg.Env,
		"port", cfg.Port,
		"log_level", cfg.LogLevel,
	)

	// Create server
	srv, err := server.New(cfg, log)
	if err != nil {
		log.Error("failed to create server", "error", err)
		os.Exit(1)
	}

	// Channel to listen for errors from server
	serverErrors := make(chan error, 1)

	// Start server in goroutine
	go func() {
		log.Info("server listening", "port", cfg.Port)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Channel to listen for interrupt signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or error
	select {
	case err := <-serverErrors:
		log.Error("server error", "error", err)
		os.Exit(1)

	case sig := <-shutdown:
		log.Info("shutdown signal received", "signal", sig.String())

		// Create context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			log.Error("graceful shutdown failed", "error", err)
			os.Exit(1)
		}

		log.Info("server stopped gracefully")
	}
}
```

> 🧠 **Config loading:**
> 📱 "config.Load() — .env file-ൽ നിന്ന് port, database path, session keys എല്ലാം load ചെയ്യും. Fail ആയാൽ panic — config ഇല്ലാതെ server start ചെയ്യണ്ട."

> 🧠 **Logger — dev vs prod:**
> 📱 "Development mode-ൽ logger.NewDevelopment() — pretty, human-readable logs. Production-ൽ logger.New() — structured JSON logs. cfg.IsDevelopment() check ചെയ്ത് switch ചെയ്യും."

> 🧠 **Server creation:**
> 📱 "server.New(cfg, log) — previous section-ൽ define ചെയ്ത Server struct create ചെയ്യും. DB connect, repositories init, middleware setup, router wire — എല്ലാം ഇതിൽ."

> 🧠 **Graceful shutdown — signal handling:**
> 📱 "Server goroutine-ൽ run ചെയ്യും. Main goroutine block ചെയ്ത് wait ചെയ്യും — `select` statement use ചെയ്ത്."
> 📱 "`os.Interrupt` — Ctrl+C press ചെയ്യുമ്പോൾ. `syscall.SIGTERM` — Docker/Kubernetes container stop ചെയ്യുമ്പോൾ."
> 📱 "Signal receive ചെയ്താൽ 10 second timeout-ൽ graceful shutdown. Active requests complete ചെയ്യാൻ time കൊടുക്കും, പക്ഷേ 10 seconds കഴിഞ്ഞാൽ force stop."

> 🧠 **`select` pattern:**
> 📱 "Go-യിലെ powerful pattern — two channels concurrent ആയി listen ചെയ്യും. Server error വന്നാൽ first case. Shutdown signal വന്നാൽ second case. ഏത് ആദ്യം വരുന്നോ അത് execute ചെയ്യും."

---

> 🎥 **Transition:** "Backend complete! ഇനി frontend — templates."

---

## 🎬 Part 2 Ending — CTA

> 🔊 **Outro music fade in**

📱 **Narration**:
> "Backend complete ആയി!"
>
> "Database, auth, CRUD, click tracking, analytics — full backend working."
>
> "Part 3-ൽ frontend build ചെയ്യും — HTMX magic, Alpine.js, beautiful Tailwind UI, dark mode, drag-drop."
>
> "ഇത് miss ആകരുത് — Subscribe ചെയ്യൂ, bell icon press ചെയ്യൂ!"
>
> "Code GitHub-ൽ ഉണ്ട്. അടുത്ത video-ൽ കാണാം!"

> 🔊 **End screen**: Subscribe button + Part 3 preview card

---

## 📝 Part 2 Editing Notes

- Add end screen with subscribe + next video link
- Quick montage of what Part 3 will look like (show finished app screenshots)
- Total Part 2 runtime target: ~45 minutes
