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

// Server holds the HTTP server and dependencies
type Server struct {
	httpServer *http.Server
	log        *slog.Logger
}

// New creates a new Server instance
func New(cfg *config.Config, log *slog.Logger) (*Server, error) {
	// Initialize database
	db, err := repository.NewDB(cfg.DatabasePath, log)
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := repository.Migrate(db, log); err != nil {
		return nil, err
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	linkRepo := repository.NewLinkRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)

	// Initialize responder
	resp := response.New(log)

	// Initialize middleware
	mw := middleware.New(log, cfg.SessionSecret, cfg.SessionEncKey, cfg.RateLimit)

	// Initialize handlers
	h := handler.New(&handler.Dependencies{
		Log:           log,
		Responder:     resp,
		Store:         mw.Store(),
		UserRepo:      userRepo,
		LinkRepo:      linkRepo,
		AnalyticsRepo: analyticsRepo,
	})

	// Initialize router
	r := router.New(h, mw)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		log:        log,
	}, nil
}

// Start begins listening for requests
func (s *Server) Start() error {
	s.log.Info("server starting", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("server shutting down")
	return s.httpServer.Shutdown(ctx)
}
