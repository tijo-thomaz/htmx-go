# Scene 10: Dashboard Handler (43:00 - 45:00)

> 🎬 **Previous**: Link handler + click tracking (Scene 9)
> 🎯 **Goal**: Dashboard with analytics summary

---

**⌨️ Create `internal/handler/dashboard.go`:**
```go
package handler

import (
	"html/template"
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

// Stats returns the stats partial HTML for HTMX polling
func (h *DashboardHandler) Stats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	analytics, err := h.analyticsRepo.GetSummary(r.Context(), userID, 28)
	if err != nil {
		h.log.Error("analytics error", "error", err)
		analytics = nil
	}

	tmpl, err := template.ParseFiles("web/templates/partials/stats.html")
	if err != nil {
		h.log.Error("template error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, analytics); err != nil {
		h.log.Error("template error", "error", err)
	}
}
```

> ⚠️ **Note:** Import `"html/template"` also needed at the top for the Stats handler.

> 🧠 **Explain:**
> 📱 "GetSummary(ctx, userID, 28) — last 28 days analytics. TotalViews, TotalClicks, per-link breakdown."
> 📱 "Analytics error ആയാൽ nil set ചെയ്യും. Dashboard render ചെയ്യും — analytics optional."
> 📱 "Template-ൽ `{{if .Analytics}}` check ചെയ്യും. nil ആയാൽ 0 show."

> 🧠 **Explain Stats handler:**
> 📱 "Stats() — HTMX polling endpoint. Dashboard-ൽ stats cards every 10 seconds auto-refresh ആകും."
> 📱 "stats.html partial template render ചെയ്ത് return ചെയ്യുന്നു — full page reload വേണ്ട."
> 📱 "Profile visit ചെയ്യുമ്പോഴോ link click ചെയ്യുമ്പോഴോ dashboard-ൽ real-time ആയി count update ആകും."

> 🎯 📱 "ഇവിടെ click tracking work ചെയ്യുന്നുണ്ടെങ്കിൽ TotalClicks > 0 കാണും. 0 ആണെങ്കിൽ Scene 9-ൽ discuss ചെയ്ത bugs check ചെയ്യൂ!"

---

## Profile Handler — Public Profile + Page View Tracking

> 📱 "ഇനി Profile handler — public profile page. ആരും login ഇല്ലാതെ `/u/{username}` visit ചെയ്യുമ്പോൾ ഈ handler ആണ് handle ചെയ്യുന്നത്."

**⌨️ Create `internal/handler/profile.go`:**
```go
package handler

import (
	"context"
	"net/http"

	"linkbio/internal/model"
	"linkbio/internal/pkg/response"
	"linkbio/internal/pkg/templates"
	"linkbio/internal/repository"

	"log/slog"

	"github.com/go-chi/chi/v5"
)

type ProfileHandler struct {
	log           *slog.Logger
	resp          *response.Responder
	userRepo      *repository.UserRepository
	linkRepo      *repository.LinkRepository
	analyticsRepo *repository.AnalyticsRepository
}

func NewProfileHandler(deps *Dependencies) *ProfileHandler {
	return &ProfileHandler{
		log:           deps.Log,
		resp:          deps.Responder,
		userRepo:      deps.UserRepo,
		linkRepo:      deps.LinkRepo,
		analyticsRepo: deps.AnalyticsRepo,
	}
}

type ProfileData struct {
	User  *model.User
	Links []model.Link
}

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

	links, err := h.linkRepo.GetActiveByUserID(r.Context(), user.ID)
	if err != nil {
		h.log.Error("database error", "error", err)
		h.resp.Error(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Record page view asynchronously
	// ⚠️ Use context.Background(), NOT r.Context()!
	// r.Context() gets cancelled after response is sent, killing the DB write.
	go func() {
		h.analyticsRepo.RecordPageView(context.Background(), user.ID, r.Referer(), r.UserAgent())
	}()

	data := ProfileData{
		User:  user,
		Links: links,
	}

	if err := templates.Render(w, "profile.html", data); err != nil {
		h.log.Error("template error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
```

> 🧠 **Explain:**
> 📱 "`chi.URLParam(r, \"username\")` — URL-ൽ നിന്ന് username extract ചെയ്യുന്നു. Route `/u/{username}` ആയതുകൊണ്ട് chi automatically parse ചെയ്യും."
> 📱 "`GetActiveByUserID` — active links മാത്രം fetch ചെയ്യും. Dashboard-ൽ `GetByUserID` ആണ് — owner-ന് inactive links-ഉം കാണണം. Public profile-ൽ inactive links show ചെയ്യരുത്."
> 📱 "`ProfileData` struct — template-ന് User-ഉം Links-ഉം pass ചെയ്യാൻ. Dashboard-ൽ Analytics-ഉം ഉണ്ട്, Profile-ൽ വേണ്ട."

> ⚠️ 📱 "Page view tracking-ൽ same pattern — `context.Background()` use ചെയ്യുന്നു. Scene 9-ൽ click tracking-ൽ discuss ചെയ്ത bug ഓർമ്മയുണ്ടോ? Same issue ഇവിടെയും. Response send ചെയ്താൽ `r.Context()` cancel ആകും, goroutine-ലെ DB write fail ആകും. `context.Background()` = independent context, cancel ആകില്ല."

> 🎯 📱 "Profile page public ആണ്. Auth middleware behind അല്ല. Router-ൽ public group-ൽ `/u/{username}` register ചെയ്യും — Scene 11-ൽ കാണാം."

---

> 🎥 **Transition:** "Handlers all done. ഇനി Router — routes wire ചെയ്യാം."
