package handler

import (
	"log/slog"

	"linkbio/internal/pkg/response"
	"linkbio/internal/repository"

	"github.com/gorilla/sessions"
)

// Handler holds all HTTP handlers
type Handler struct {
	Auth      *AuthHandler
	Link      *LinkHandler
	Profile   *ProfileHandler
	Dashboard *DashboardHandler
	Health    *HealthHandler
}

// Dependencies for handlers
type Dependencies struct {
	Log         *slog.Logger
	Responder   *response.Responder
	Store       *sessions.CookieStore
	UserRepo    *repository.UserRepository
	LinkRepo    *repository.LinkRepository
	AnalyticsRepo *repository.AnalyticsRepository
}

// New creates all handlers
func New(deps *Dependencies) *Handler {
	return &Handler{
		Auth:      NewAuthHandler(deps),
		Link:      NewLinkHandler(deps),
		Profile:   NewProfileHandler(deps),
		Dashboard: NewDashboardHandler(deps),
		Health:    NewHealthHandler(deps.Log),
	}
}
