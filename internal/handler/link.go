package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"linkbio/internal/middleware"
	"linkbio/internal/model"
	"linkbio/internal/pkg/response"
	"linkbio/internal/repository"

	"log/slog"

	"github.com/go-chi/chi/v5"
)

// LinkHandler handles link CRUD endpoints
type LinkHandler struct {
	log      *slog.Logger
	resp     *response.Responder
	linkRepo *repository.LinkRepository
	analyticsRepo *repository.AnalyticsRepository
}

// NewLinkHandler creates a new LinkHandler
func NewLinkHandler(deps *Dependencies) *LinkHandler {
	return &LinkHandler{
		log:      deps.Log,
		resp:     deps.Responder,
		linkRepo: deps.LinkRepo,
		analyticsRepo: deps.AnalyticsRepo,
	}
}

// Create adds a new link
func (h *LinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == 0 {
		h.resp.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := r.ParseForm(); err != nil {
		h.resp.Error(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	title := r.FormValue("title")
	url := r.FormValue("url")

	if title == "" || url == "" {
		h.resp.Error(w, http.StatusBadRequest, "Title and URL are required")
		return
	}

	link := &model.Link{
		UserID:   userID,
		Title:    title,
		URL:      url,
		Icon:     r.FormValue("icon"),
		IsActive: true,
	}

	if err := h.linkRepo.Create(r.Context(), link); err != nil {
		h.log.Error("link creation error", "error", err)
		h.resp.Error(w, http.StatusInternalServerError, "Failed to create link")
		return
	}

	h.log.Info("link created", "link_id", link.ID, "user_id", userID)

	// Return the new link as HTML partial for HTMX
	tmpl, err := template.ParseFiles("web/templates/partials/link.html")
	if err != nil {
		h.log.Error("template error", "error", err)
		h.resp.Error(w, http.StatusInternalServerError, "Template error")
		return
	}
	tmpl.Execute(w, link)

	// OOB: update link count badge
	count, _ := h.linkRepo.CountByUserID(r.Context(), userID)
	fmt.Fprintf(w, `<span id="link-count" hx-swap-oob="true" class="ml-2 px-2 py-0.5 text-xs font-medium rounded-full bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400">%d</span>`, count)

	// OOB: hide empty state
	fmt.Fprint(w, `<div id="empty-state" hx-swap-oob="outerHTML" style="display:none"></div>`)
}

// Update modifies an existing link
func (h *LinkHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == 0 {
		h.resp.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	linkID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.resp.Error(w, http.StatusBadRequest, "Invalid link ID")
		return
	}

	// Get existing link
	link, err := h.linkRepo.GetByID(r.Context(), linkID)
	if err != nil {
		h.log.Error("database error", "error", err)
		h.resp.Error(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if link == nil || link.UserID != userID {
		h.resp.Error(w, http.StatusNotFound, "Link not found")
		return
	}

	if err := r.ParseForm(); err != nil {
		h.resp.Error(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	link.Title = r.FormValue("title")
	link.URL = r.FormValue("url")
	link.Icon = r.FormValue("icon")
	link.IsActive = r.FormValue("is_active") == "on" || r.FormValue("is_active") == "true"

	if err := h.linkRepo.Update(r.Context(), link); err != nil {
		h.log.Error("link update error", "error", err)
		h.resp.Error(w, http.StatusInternalServerError, "Failed to update link")
		return
	}

	h.log.Info("link updated", "link_id", link.ID, "user_id", userID)

	tmpl, _ := template.ParseFiles("web/templates/partials/link.html")
	tmpl.Execute(w, link)
}

// Delete removes a link
func (h *LinkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == 0 {
		h.resp.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	linkID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.resp.Error(w, http.StatusBadRequest, "Invalid link ID")
		return
	}

	// Verify ownership
	link, err := h.linkRepo.GetByID(r.Context(), linkID)
	if err != nil || link == nil || link.UserID != userID {
		h.resp.Error(w, http.StatusNotFound, "Link not found")
		return
	}

	if err := h.linkRepo.Delete(r.Context(), linkID); err != nil {
		h.log.Error("link delete error", "error", err)
		h.resp.Error(w, http.StatusInternalServerError, "Failed to delete link")
		return
	}

	h.log.Info("link deleted", "link_id", linkID, "user_id", userID)

	// OOB: update link count badge
	count, _ := h.linkRepo.CountByUserID(r.Context(), userID)
	fmt.Fprintf(w, `<span id="link-count" hx-swap-oob="true" class="ml-2 px-2 py-0.5 text-xs font-medium rounded-full bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400">%d</span>`, count)

	// OOB: restore empty state if no links remain
	if count == 0 {
		fmt.Fprint(w, `<div id="empty-state" hx-swap-oob="outerHTML" class="p-12 text-center"><div class="w-16 h-16 mx-auto mb-4 rounded-2xl bg-gray-100 dark:bg-gray-800 flex items-center justify-center"><svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/></svg></div><h3 class="text-lg font-medium text-gray-900 dark:text-white mb-1">No links yet</h3><p class="text-gray-500 dark:text-gray-400">Add your first link to get started</p></div>`)
	}
}

// Reorder updates link positions
func (h *LinkHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == 0 {
		h.resp.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var positions map[int64]int
	if err := json.NewDecoder(r.Body).Decode(&positions); err != nil {
		h.resp.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.linkRepo.UpdatePositions(r.Context(), userID, positions); err != nil {
		h.log.Error("reorder error", "error", err)
		h.resp.Error(w, http.StatusInternalServerError, "Failed to reorder links")
		return
	}

	h.log.Info("links reordered", "user_id", userID)

	w.WriteHeader(http.StatusOK)
}

// Click records a link click
func (h *LinkHandler) Click(w http.ResponseWriter, r *http.Request) {
	linkID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.resp.Error(w, http.StatusBadRequest, "Invalid link ID")
		return
	}

	link, err := h.linkRepo.GetByID(r.Context(), linkID)
	if err != nil || link == nil {
		h.resp.Error(w, http.StatusNotFound, "Link not found")
		return
	}

	go func() {
		if err := h.analyticsRepo.RecordLinkClick(context.Background(), link.UserID, linkID, r.Referer(), r.UserAgent()); err != nil {
			h.log.Error("failed to record click", "link_id", linkID, "error", err)
		}
	}()

	// Redirect to the actual URL
	http.Redirect(w, r, link.URL, http.StatusTemporaryRedirect)
}
