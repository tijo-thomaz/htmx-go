package handler

import (
	"net/http"

	"linkbio/internal/model"
	"linkbio/internal/pkg/response"
	"linkbio/internal/pkg/templates"
	"linkbio/internal/repository"

	"log/slog"

	"github.com/go-chi/chi/v5"
)

// ProfileHandler handles public profile endpoints
type ProfileHandler struct {
	log           *slog.Logger
	resp          *response.Responder
	userRepo      *repository.UserRepository
	linkRepo      *repository.LinkRepository
	analyticsRepo *repository.AnalyticsRepository
}

// NewProfileHandler creates a new ProfileHandler
func NewProfileHandler(deps *Dependencies) *ProfileHandler {
	return &ProfileHandler{
		log:           deps.Log,
		resp:          deps.Responder,
		userRepo:      deps.UserRepo,
		linkRepo:      deps.LinkRepo,
		analyticsRepo: deps.AnalyticsRepo,
	}
}

// ProfileData holds data for the profile template
type ProfileData struct {
	User  *model.User
	Links []model.Link
}

// Show renders a user's public profile
func (h *ProfileHandler) Show(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	user, err := h.userRepo.GetByUsername(r.Context(), username)
	if err != nil {
		h.log.Error("database error", "error", err)
		h.resp.Error(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if user == nil {
		h.resp.Error(w, http.StatusNotFound, "Profile not found")
		return
	}

	// Get active links
	links, err := h.linkRepo.GetActiveByUserID(r.Context(), user.ID)
	if err != nil {
		h.log.Error("database error", "error", err)
		h.resp.Error(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Record page view asynchronously
	go func() {
		h.analyticsRepo.RecordPageView(r.Context(), user.ID, r.Referer(), r.UserAgent())
	}()

	h.log.Info("profile data", "username", username, "user_id", user.ID, "links_count", len(links))
	for i, l := range links {
		h.log.Info("link", "index", i, "id", l.ID, "title", l.Title, "active", l.IsActive)
	}

	data := ProfileData{
		User:  user,
		Links: links,
	}

	if err := templates.Render(w, "profile.html", data); err != nil {
		h.log.Error("template error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
