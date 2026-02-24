# Scene 9: Link Handler & Click Tracking ⚠️ (38:00 - 43:00)

> 🎬 **Previous**: Auth handlers (Scene 8)
> 🎯 **Goal**: Link CRUD + click tracking with TWO production bugs avoided
> ⚠️ **Teaching highlight**: This scene has the two biggest "aha moments" of the tutorial

---

## 🎥 Camera Note

> 📱 "ഈ section-ൽ Link CRUD-ഉം click tracking-ഉം build ചെയ്യും. ⚠️ Click tracking-ൽ production apps-ൽ commonly കാണുന്ന രണ്ട് bugs avoid ചെയ്യും. ശ്രദ്ധിക്കൂ!"

---

## Link Handler — Full Code

**⌨️ Create `internal/handler/link.go`:**
```go
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

type LinkHandler struct {
	log           *slog.Logger
	resp          *response.Responder
	linkRepo      *repository.LinkRepository
	analyticsRepo *repository.AnalyticsRepository
}

func NewLinkHandler(deps *Dependencies) *LinkHandler {
	return &LinkHandler{
		log:           deps.Log,
		resp:          deps.Responder,
		linkRepo:      deps.LinkRepo,
		analyticsRepo: deps.AnalyticsRepo,
	}
}
```

> 🧠 📱 "LinkHandler-ന് analyticsRepo-ഉം inject ചെയ്യുന്നു — click tracking-ന്."
> 📱 "import `\"context\"` — ⚠️ ഇത് important. Click handler-ൽ context.Background() use ചെയ്യും. എന്തുകൊണ്ടെന്ന് soon explain ചെയ്യാം."

---

## Create Link — with HTMX OOB Swaps

**⌨️ Continue in link.go:**
```go
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
		UserID: userID, Title: title, URL: url,
		Icon: r.FormValue("icon"), IsActive: true,
	}

	if err := h.linkRepo.Create(r.Context(), link); err != nil {
		h.log.Error("link creation error", "error", err)
		h.resp.Error(w, http.StatusInternalServerError, "Failed to create link")
		return
	}

	h.log.Info("link created", "link_id", link.ID, "user_id", userID)

	// Return new link HTML for HTMX
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
```

> 🧠 **Explain OOB (Out-of-Band) swaps:**
> 📱 "HTMX response-ൽ main content + OOB elements ഒരുമിച്ച് send ചെയ്യാം."
> 📱 "Main: new link HTML → #links-list-ൽ insert."
> 📱 "OOB 1: link count badge update. hx-swap-oob='true' — matching id find ചെയ്ത് replace."
> 📱 "OOB 2: empty state hide. Links ഉണ്ടെങ്കിൽ 'No links yet' message hide ചെയ്യണം."
> 📱 "ഒരു response, മൂന്ന് updates. Page reload ഇല്ല!"

---

## Link Partial Template

> 📱 "Create handler return ചെയ്യുന്ന HTML template ഇതാണ്. HTMX response-ൽ ഈ partial render ചെയ്ത് send ചെയ്യും."

**⌨️ Create `web/templates/partials/link.html`:**
```html
<div class="link-card flex items-center gap-4 p-5 hover:bg-gray-50 dark:hover:bg-gray-800/50" 
     data-link-id="{{.ID}}">
    <button class="drag-handle cursor-grab active:cursor-grabbing p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8h16M4 16h16"/>
        </svg>
    </button>
    <div class="flex-1 min-w-0">
        <h3 class="font-medium text-gray-900 dark:text-white truncate">{{.Title}}</h3>
        <p class="text-sm text-gray-500 dark:text-gray-400 truncate">{{.URL}}</p>
    </div>
    <button hx-delete="/api/v1/links/{{.ID}}"
            hx-target="closest .link-card"
            hx-swap="outerHTML swap:200ms"
            hx-confirm="Delete this link?"
            class="p-2 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
        </svg>
    </button>
</div>
```

> 🧠 **Explain template structure:**
> 📱 "ഈ partial template-ആണ് Create handler-ൽ `template.ParseFiles(\"web/templates/partials/link.html\")` load ചെയ്യുന്നത്. HTMX response-ൽ ഈ HTML client-ന് send ചെയ്യും."
> 📱 "`data-link-id=\"{{.ID}}\"` — SortableJS drag-drop reorder-ന് ഈ attribute use ചെയ്യും. ഓരോ link-ന്റെയും ID track ചെയ്യാൻ."
> 📱 "`drag-handle` class — ഈ button മാത്രം drag trigger ചെയ്യും. Card-ൽ anywhere drag ചെയ്യാൻ allow ചെയ്യില്ല."
> 📱 "`hx-delete=\"/api/v1/links/{{.ID}}\"` — delete button click ചെയ്യുമ്പോൾ HTMX DELETE request send ചെയ്യും."
> 📱 "`hx-target=\"closest .link-card\"` — parent `.link-card` div find ചെയ്യും. ഈ entire card-ആണ് target."
> 📱 "`hx-swap=\"outerHTML swap:200ms\"` — card full remove ചെയ്യും, 200ms fade animation-ഓടെ. Smooth UX."
> 📱 "`hx-confirm=\"Delete this link?\"` — browser confirmation dialog show ചെയ്യും. Accidental delete prevent ചെയ്യാൻ."
> 📱 "`truncate` class — long titles-ഉം URLs-ഉം ellipsis (...) ആയി cut ചെയ്യും. Layout break ആകില്ല."

---

## Update, Delete, Reorder

**⌨️ Continue:**
```go
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

	tmpl, _ := template.ParseFiles("web/templates/partials/link.html")
	tmpl.Execute(w, link)
}

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

	// OOB: update count + restore empty state if needed
	count, _ := h.linkRepo.CountByUserID(r.Context(), userID)
	fmt.Fprintf(w, `<span id="link-count" hx-swap-oob="true" class="ml-2 px-2 py-0.5 text-xs font-medium rounded-full bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400">%d</span>`, count)

	if count == 0 {
		fmt.Fprint(w, `<div id="empty-state" hx-swap-oob="outerHTML" class="p-12 text-center"><div class="w-16 h-16 mx-auto mb-4 rounded-2xl bg-gray-100 dark:bg-gray-800 flex items-center justify-center"><svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/></svg></div><h3 class="text-lg font-medium text-gray-900 dark:text-white mb-1">No links yet</h3><p class="text-gray-500 dark:text-gray-400">Add your first link to get started</p></div>`)
	}
}

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
	w.WriteHeader(http.StatusOK)
}
```

> 🧠 📱 "Update — ownership check: link.UserID != userID. Other users' links edit ചെയ്യാൻ പറ്റില്ല."
> 📱 "Delete — count 0 ആയാൽ empty state HTML OOB restore ചെയ്യും."
> 📱 "Reorder — JSON body decode, transaction-based position update."

---

## ⚠️ Click Handler — THE BIG ONE

> 🎥 **Camera note:** Slow down here. This is the most important teaching moment.

**⌨️ Continue:**
```go
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

	// ⚠️ BUG FIX: context.Background(), NOT r.Context()!
	go func() {
		if err := h.analyticsRepo.RecordLinkClick(context.Background(), link.UserID, linkID, r.Referer(), r.UserAgent()); err != nil {
			h.log.Error("failed to record click", "link_id", linkID, "error", err)
		}
	}()

	http.Redirect(w, r, link.URL, http.StatusTemporaryRedirect)
}
```

---

### ⚠️ Bug 1: Context Cancellation in Goroutines

> 📱 "ഈ line നോക്കൂ: `context.Background()`. `r.Context()` use ചെയ്യാത്തത് എന്തുകൊണ്ട്?"

> 📱 "http.Redirect() call ചെയ്യുമ്പോൾ response user-ന് send ചെയ്യും. Go automatically r.Context() cancel ചെയ്യും."
> 📱 "But goroutine still running ആണ്! r.Context() already cancelled — database write silently fail! Error log-ൽ പോലും കാണില്ല (unless we check)."
> 📱 "context.Background() = independent context. Cancel ആകില്ല. Goroutine complete ചെയ്യും."

> 🎯 **Analogy:**
> 📱 "Restaurant-ൽ order കൊടുത്ത് walk out ചെയ്യുന്നത് imagine ചെയ്യൂ."
> 📱 "r.Context() = 'Customer left, order cancel ചെയ്യൂ'. Kitchen work stop."
> 📱 "context.Background() = 'Customer left, but order complete ചെയ്യൂ, takeaway pack ചെയ്യൂ'. Kitchen continue."

> ⚠️ 📱 "ഇത് very common bug ആണ്. Production apps-ൽ analytics data missing ആകുന്നത് ഇതുകൊണ്ടാണ്. Dashboard-ൽ 0 clicks."

---

### ⚠️ Bug 2: This Route Must Be Public

> 📱 "ഈ Click endpoint auth middleware-ന്റെ behind ഇടരുത്! എന്തുകൊണ്ട്?"
> 📱 "Profile page `/u/username` — public page. Visitors logged in അല്ല."
> 📱 "Visitor link click ചെയ്യുമ്പോൾ `/click/5` hit ചെയ്യും."
> 📱 "Auth middleware ഉണ്ടെങ്കിൽ → login page redirect → click tracking fail!"
> 📱 "ഇത് Router scene-ൽ (Scene 11) handle ചെയ്യും. `/click/{id}` public group-ൽ ഇടും."

> 🎯 📱 "ഒരു line code mistake — dashboard-ൽ permanently 0 clicks. Error message ഇല്ല. Silently broken. ഇതാണ് production bugs — loud crash അല്ല, silent failure."

---

> 🎥 **Transition:** "Link handler done. ഇനി Dashboard handler."
