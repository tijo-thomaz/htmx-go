package handler

import (
	"net/http"

	"linkbio/internal/middleware"
	"linkbio/internal/model"
	"linkbio/internal/pkg/response"
	"linkbio/internal/pkg/templates"
	"linkbio/internal/repository"

	"log/slog"
)

// DashboardHandler handles dashboard endpoints
type DashboardHandler struct {
	log           *slog.Logger
	resp          *response.Responder
	userRepo      *repository.UserRepository
	linkRepo      *repository.LinkRepository
	analyticsRepo *repository.AnalyticsRepository
}

// NewDashboardHandler creates a new DashboardHandler
func NewDashboardHandler(deps *Dependencies) *DashboardHandler {
	return &DashboardHandler{
		log:           deps.Log,
		resp:          deps.Responder,
		userRepo:      deps.UserRepo,
		linkRepo:      deps.LinkRepo,
		analyticsRepo: deps.AnalyticsRepo,
	}
}

// DashboardData holds data for the dashboard template
type DashboardData struct {
	User      *model.User
	Links     []model.Link
	Analytics *model.AnalyticsSummary
}

// Index renders the dashboard
func (h *DashboardHandler) Index(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	username := middleware.UsernameFromContext(r.Context())

	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		h.log.Error("database error", "error", err)
		h.resp.Error(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	links, err := h.linkRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		h.log.Error("database error", "error", err)
		h.resp.Error(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	analytics, err := h.analyticsRepo.GetSummary(r.Context(), userID, 28) // Last 28 days
	if err != nil {
		h.log.Error("analytics error", "error", err)
		// Continue without analytics
		analytics = nil
	}

	h.log.Debug("dashboard loaded", "user_id", userID, "username", username, "links_count", len(links))

	data := DashboardData{
		User:      user,
		Links:     links,
		Analytics: analytics,
	}

	if err := templates.Render(w, "dashboard.html", data); err != nil {
		h.log.Error("template error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
