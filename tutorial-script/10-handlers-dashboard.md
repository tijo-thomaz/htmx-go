# Scene 10: Dashboard Handler (43:00 - 45:00)

> 🎬 **Previous**: Link handler + click tracking (Scene 9)
> 🎯 **Goal**: Dashboard with analytics summary

---

**⌨️ Create `internal/handler/dashboard.go`:**
```go
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

type DashboardHandler struct {
	log           *slog.Logger
	resp          *response.Responder
	userRepo      *repository.UserRepository
	linkRepo      *repository.LinkRepository
	analyticsRepo *repository.AnalyticsRepository
}

func NewDashboardHandler(deps *Dependencies) *DashboardHandler {
	return &DashboardHandler{
		log:           deps.Log,
		resp:          deps.Responder,
		userRepo:      deps.UserRepo,
		linkRepo:      deps.LinkRepo,
		analyticsRepo: deps.AnalyticsRepo,
	}
}

type DashboardData struct {
	User      *model.User
	Links     []model.Link
	Analytics *model.AnalyticsSummary
}

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

	analytics, err := h.analyticsRepo.GetSummary(r.Context(), userID, 28)
	if err != nil {
		h.log.Error("analytics error", "error", err)
		analytics = nil // Continue without analytics
	}

	h.log.Debug("dashboard loaded", "user_id", userID, "username", username, "links_count", len(links))

	data := DashboardData{User: user, Links: links, Analytics: analytics}

	if err := templates.Render(w, "dashboard.html", data); err != nil {
		h.log.Error("template error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
```

> 🧠 **Explain:**
> 📱 "GetSummary(ctx, userID, 28) — last 28 days analytics. TotalViews, TotalClicks, per-link breakdown."
> 📱 "Analytics error ആയാൽ nil set ചെയ്യും. Dashboard render ചെയ്യും — analytics optional."
> 📱 "Template-ൽ `{{if .Analytics}}` check ചെയ്യും. nil ആയാൽ 0 show."

> 🎯 📱 "ഇവിടെ click tracking work ചെയ്യുന്നുണ്ടെങ്കിൽ TotalClicks > 0 കാണും. 0 ആണെങ്കിൽ Scene 9-ൽ discuss ചെയ്ത bugs check ചെയ്യൂ!"

---

> 🎥 **Transition:** "Handlers all done. ഇനി Router — routes wire ചെയ്യാം."
