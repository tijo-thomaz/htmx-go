# 🎬 PART 2 — Backend & Logic

> 📱 **Part 2 Intro** (show on screen before Scene 5):
> "Part 1-ൽ നമ്മൾ project setup, config, logging — foundation build ചെയ്തു. Server running ആണ്."
>
> "ഇന്ന് real code — database, authentication, handlers — full backend build ചെയ്യും!"
>
> "Ready ആണോ? Let's go!"

> 🔊 **Intro music** (3 seconds)

---

# Scene 5: Database Models (20:00 - 23:00)

> 🎬 **Previous Scene**: Logging (Scene 4) — we built structured logging with `slog`
> 🎬 **Camera**: VS Code open, terminal visible. Switch to `internal/model/` folder in the file explorer.
> 🎬 **Transition**: "Logging ready ആയി. ഇനി നമ്മുടെ app-ന്റെ data define ചെയ്യാം."

---

## Opening

**📱 Read this:**
> "ആദ്യം data structures define ചെയ്യാം. ഇത് database tables-ന്റെ Go representation ആണ്. Database-ൽ എന്തൊക്കെ store ചെയ്യണം എന്ന് ഇവിടെ decide ചെയ്യും."

---

## User Model

**📱 Read this:**
> "ആദ്യത്തെ model — User. Register ചെയ്യുന്ന ഓരോ user-ന്റെയും details ഇതിൽ store ചെയ്യും."

**⌨️ Create `internal/model/user.go`:**
```go
package model

import "time"

// User represents a registered user
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	DisplayName  string    `json:"display_name"`
	Bio          string    `json:"bio"`
	AvatarURL    string    `json:"avatar_url"`
	Theme        string    `json:"theme"`
	CreatedAt    time.Time `json:"created_at"`
}
```

**🧠 Explain:**
> "struct tags നോക്കൂ — `json:"id"`, `json:"username"` — ഇത് Go struct-നെ JSON-ലേക്ക് convert ചെയ്യുമ്പോൾ field name define ചെയ്യുന്നു."

**⚠️ Security — `json:"-"`:**
> "PasswordHash-ന്റെ tag `json:\"-\"` ആണ്. ഇതിന്റെ meaning: ഈ field ഒരിക്കലും JSON response-ൽ include ചെയ്യില്ല. Password hash client-ലേക്ക് അയക്കരുത് — അത് serious security risk ആണ്!"

---

## Link Model

**📱 Read this:**
> "രണ്ടാമത്തെ model — Link. User-ന്റെ profile page-ൽ കാണിക്കുന്ന ഓരോ link-ഉം ഇതിൽ store ചെയ്യും."

**⌨️ Create `internal/model/link.go`:**
```go
package model

import "time"

// Link represents a user's link
type Link struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Icon      string    `json:"icon"`
	Position  int       `json:"position"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
```

**🧠 Explain `Position`:**
> "Position field — ഇത് link-ന്റെ order define ചെയ്യുന്നു. User dashboard-ൽ drag-and-drop ചെയ്ത് links reorder ചെയ്യുമ്പോൾ ഈ value update ആകും. Position 1 ആണ് top-ൽ, 2 രണ്ടാമത്, അങ്ങനെ."

**🧠 Explain `IsActive`:**
> "IsActive — ഒരു toggle switch പോലെ. User-ന് ഒരു link delete ചെയ്യാതെ temporarily hide ചെയ്യാൻ പറ്റും. IsActive false ആണെങ്കിൽ public profile-ൽ ആ link കാണില്ല."

---

## Analytics Model

**📱 Read this:**
> "മൂന്നാമത്തെ model — Analytics. ഇത് important ആണ്. ആര് profile visit ചെയ്തു, ഏത് link click ചെയ്തു — എല്ലാം track ചെയ്യും."

**⌨️ Create `internal/model/analytics.go`:**
```go
package model

import "time"

// Analytics represents a tracking event
type Analytics struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	LinkID    *int64    `json:"link_id,omitempty"` // nil for page views
	EventType string    `json:"event_type"`        // "page_view" or "link_click"
	Referrer  string    `json:"referrer"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// AnalyticsSummary holds aggregated analytics data
type AnalyticsSummary struct {
	TotalViews  int              `json:"total_views"`
	TotalClicks int              `json:"total_clicks"`
	LinkClicks  []LinkClickCount `json:"link_clicks"`
}

// LinkClickCount holds click count for a specific link
type LinkClickCount struct {
	LinkID int64  `json:"link_id"`
	Title  string `json:"title"`
	Clicks int    `json:"clicks"`
}
```

**🧠 Explain — Why `LinkID` is `*int64` (pointer):**
> "ഇവിടെ ഒരു important concept ഉണ്ട്. LinkID-യുടെ type `*int64` ആണ് — star int64, ഒരു pointer. എന്താ reason?"
>
> "നമ്മുടെ app-ൽ രണ്ട് type events ഉണ്ട്: `page_view` ഉം `link_click` ഉം. ആരെങ്കിലും profile page visit ചെയ്താൽ — അത് page_view. അവിടെ specific link ഇല്ല, so LinkID nil ആകും. ആരെങ്കിലും ഒരു link click ചെയ്താൽ — അത് link_click. അപ്പോൾ LinkID-ൽ ആ link-ന്റെ ID ഉണ്ടാകും."
>
> "Regular `int64` use ചെയ്താൽ nil ആക്കാൻ പറ്റില്ല — default 0 ആകും. 0 ഒരു valid ID ആണോ അല്ലയോ confuse ആകും. Pointer use ചെയ്താൽ nil means 'no link', value means 'this link' — clear!"

**🧠 Explain `AnalyticsSummary`:**
> "AnalyticsSummary — ഇത് dashboard-ൽ കാണിക്കാൻ ഉള്ള aggregated data. Total views എത്ര, total clicks എത്ര — ഒറ്റ query-ൽ കിട്ടും."

**🧠 Explain `LinkClickCount`:**
> "LinkClickCount — ഏത് link-ന് ആണ് കൂടുതൽ clicks എന്ന് per-link breakdown കാണിക്കും. Dashboard-ൽ ഒരു table ആയി display ചെയ്യാം — 'Instagram: 45 clicks, YouTube: 30 clicks' എന്ന format-ൽ."

---

## Analogy

**🎯 Analogy:**
> "ഒരു simple analogy പറയാം. Models ഒരു house-ന്റെ blueprint പോലെ ആണ്. Blueprint-ൽ define ചെയ്യും — ഇവിടെ bedroom, ഇവിടെ kitchen, ഇത്ര size. Database ആ blueprint use ചെയ്ത് build ചെയ്ത actual house ആണ്. Blueprint ഇല്ലാതെ house build ചെയ്യാൻ പറ്റില്ല — models ഇല്ലാതെ database design ചെയ്യാൻ പറ്റില്ല."

---

## Recap

**📱 Read this:**
> "So മൂന്ന് models ready ആയി:
> - **User** — account details, profile info
> - **Link** — user-ന്റെ links with ordering and toggle
> - **Analytics** — page views, link clicks, per-link breakdown
>
> ഇനി next step — ഈ models database-ൽ actual tables ആയി create ചെയ്യാം. Database connection and migrations."

---

> 🎬 **Next Scene**: Database Connection & Migrations (Scene 6)
> 🎬 **Transition**: Move to `internal/repository/db.go`
