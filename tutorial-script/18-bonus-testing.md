# Scene 18: Testing & Benchmarks (Bonus)

> 🎬 **Bonus section** — can be a separate video or appendix
> 🎯 **Goal**: Unit tests, integration tests, benchmarks, coverage

---

## Why Testing?

> 📱 "Production-ൽ bug കണ്ടെത്തുന്നതിനേക്കാൾ development-ൽ കണ്ടെത്തുന്നത് നല്ലത്."

> 🎯 📱 "Car-ന്റെ brakes test ചെയ്യാതെ road-ൽ ഇറക്കുമോ?"

---

## Test Utilities

**⌨️ Create `internal/testutil/testutil.go`:**
```go
package testutil

import (
	"database/sql"
	"log/slog"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

// TestDB creates an in-memory database for testing
// Uses single connection to ensure :memory: DB consistency
func TestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// IMPORTANT: SQLite :memory: creates separate DB per connection
	// Limit to 1 connection to keep tests stable
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Enable foreign keys (SQLite disables by default)
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	// Run migrations inline to avoid import cycle
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			display_name TEXT DEFAULT '',
			bio TEXT DEFAULT '',
			avatar_url TEXT DEFAULT '',
			theme TEXT DEFAULT 'light',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS links (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			url TEXT NOT NULL,
			icon TEXT DEFAULT '',
			position INTEGER DEFAULT 0,
			is_active INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS analytics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			link_id INTEGER,
			event_type TEXT NOT NULL,
			referrer TEXT DEFAULT '',
			user_agent TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (link_id) REFERENCES links(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_links_user_id ON links(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_analytics_user_id ON analytics(user_id)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			t.Fatalf("failed to run migration: %v", err)
		}
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

// TestLogger creates a silent logger for testing
func TestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}
```

> 🧠 📱 "t.Helper() — test fail ആയാൽ line number ശരിയായ place-ൽ point ചെയ്യും. Helper function-ൽ ആയാൽ confusing line number കാണും."
> 📱 ":memory: — in-memory SQLite. Test-ന് ശേഷം auto-delete. Disk I/O ഇല്ല, fast."
> 📱 "t.Cleanup() — test end-ൽ db.Close() auto-call. defer-നേക്കാൾ better — sub-tests-ൽ correct timing guarantee."
> 📱 "⚠️ Import cycle avoid ചെയ്യാൻ migrations inline ആണ്. testutil package-ൽ repository import ചെയ്താൽ cycle ആകും!"
> 📱 "SetMaxOpenConns(1) — SQLite :memory: database ഓരോ connection-നും separate ആണ്. One connection = consistent data."
> 📱 "PRAGMA foreign_keys = ON — SQLite default-ൽ foreign keys OFF ആണ്! Tests-ൽ ON ചെയ്യണം, production code-ൽ db.go-ൽ pragma ആയി ON ആണ്."
> 📱 "TestLogger() — error level only. Tests-ൽ noise reduce ചെയ്യാൻ."

---

## User Repository Tests

**⌨️ Create `internal/repository/user_test.go`:**
```go
package repository

import (
	"context"
	"testing"

	"linkbio/internal/model"
	"linkbio/internal/testutil"
)

func TestUserRepository_Create(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		DisplayName:  "Test User",
		Theme:        "light",
	}

	// Test: Create user
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Verify: ID was set
	if user.ID == 0 {
		t.Error("Create() did not set user ID")
	}

	// Verify: User exists in database
	found, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if found == nil {
		t.Fatal("GetByID() returned nil")
	}
	if found.Username != user.Username {
		t.Errorf("Username = %v, want %v", found.Username, user.Username)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Setup: Create a user
	user := &model.User{
		Username:     "emailtest",
		Email:        "find@example.com",
		PasswordHash: "hash",
		DisplayName:  "Email Test",
		Theme:        "dark",
	}
	repo.Create(ctx, user)

	// Test cases
	tests := []struct {
		name      string
		email     string
		wantFound bool
	}{
		{"existing email", "find@example.com", true},
		{"non-existing email", "notfound@example.com", false},
		{"empty email", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := repo.GetByEmail(ctx, tt.email)
			if err != nil {
				t.Fatalf("GetByEmail() error = %v", err)
			}
			gotFound := found != nil
			if gotFound != tt.wantFound {
				t.Errorf("GetByEmail() found = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func TestUserRepository_GetByUsername(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username:     "findbyname",
		Email:        "findbyname@example.com",
		PasswordHash: "hash",
		DisplayName:  "Find By Name",
		Theme:        "light",
	}
	repo.Create(ctx, user)

	tests := []struct {
		name      string
		username  string
		wantFound bool
	}{
		{"existing username", "findbyname", true},
		{"non-existing username", "notfound", false},
		{"case sensitive", "FINDBYNAME", false}, // SQLite is case-sensitive by default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := repo.GetByUsername(ctx, tt.username)
			if err != nil {
				t.Fatalf("GetByUsername() error = %v", err)
			}
			gotFound := found != nil
			if gotFound != tt.wantFound {
				t.Errorf("GetByUsername() found = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func TestUserRepository_DuplicateUsername(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user1 := &model.User{
		Username: "duplicate", Email: "user1@example.com",
		PasswordHash: "hash", DisplayName: "User 1", Theme: "light",
	}
	if err := repo.Create(ctx, user1); err != nil {
		t.Fatalf("Create() first user error = %v", err)
	}

	user2 := &model.User{
		Username: "duplicate", Email: "user2@example.com",
		PasswordHash: "hash", DisplayName: "User 2", Theme: "light",
	}
	err := repo.Create(ctx, user2)
	if err == nil {
		t.Error("Create() should fail for duplicate username")
	}
}

func TestUserRepository_DuplicateEmail(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user1 := &model.User{
		Username: "user1", Email: "duplicate@example.com",
		PasswordHash: "hash", DisplayName: "User 1", Theme: "light",
	}
	if err := repo.Create(ctx, user1); err != nil {
		t.Fatalf("Create() first user error = %v", err)
	}

	user2 := &model.User{
		Username: "user2", Email: "duplicate@example.com",
		PasswordHash: "hash", DisplayName: "User 2", Theme: "light",
	}
	err := repo.Create(ctx, user2)
	if err == nil {
		t.Error("Create() should fail for duplicate email")
	}
}
```

> 🧠 📱 "Table-driven tests — tests := []struct{}. Multiple cases, one function. Go community standard."
> 📱 "t.Run() — sub-tests. ഓരോ case independently run. Output-ൽ `TestUserRepository_GetByEmail/existing_email` format-ൽ കാണും."
> 📱 "DuplicateUsername, DuplicateEmail — database UNIQUE constraint work ചെയ്യുന്നുണ്ട് എന്ന് verify."
> 📱 "GetByUsername case sensitivity test — SQLite default-ൽ case-sensitive. 'FINDBYNAME' != 'findbyname'."

---

## Link Repository Tests

**⌨️ Create `internal/repository/link_test.go`:**
```go
package repository

import (
	"context"
	"testing"

	"linkbio/internal/model"
	"linkbio/internal/testutil"
)

// createTestUser is a helper to reduce boilerplate
func createTestUser(t *testing.T, repo *UserRepository, username string) *model.User {
	t.Helper()
	user := &model.User{
		Username:     username,
		Email:        username + "@test.com",
		PasswordHash: "hash",
		DisplayName:  username,
		Theme:        "light",
	}
	if err := repo.Create(context.Background(), user); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	return user
}

func TestLinkRepository_Create(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	user := createTestUser(t, userRepo, "linktest")

	link := &model.Link{
		UserID: user.ID, Title: "My Website",
		URL: "https://example.com", IsActive: true,
	}

	err := linkRepo.Create(ctx, link)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if link.ID == 0 {
		t.Error("Create() did not set link ID")
	}
	// Verify position was auto-set
	if link.Position != 1 {
		t.Errorf("Position = %v, want 1", link.Position)
	}
}

func TestLinkRepository_PositionAutoIncrement(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	user := createTestUser(t, userRepo, "postest")

	// Create 3 links — positions should auto-increment
	for i := 1; i <= 3; i++ {
		link := &model.Link{
			UserID: user.ID, Title: "Link",
			URL: "https://example.com", IsActive: true,
		}
		linkRepo.Create(ctx, link)
		if link.Position != i {
			t.Errorf("Link %d position = %v, want %v", i, link.Position, i)
		}
	}
}

func TestLinkRepository_GetByID(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	user := createTestUser(t, userRepo, "getbyidtest")

	link := &model.Link{
		UserID: user.ID, Title: "Test Link",
		URL: "https://test.com", Icon: "🔗", IsActive: true,
	}
	linkRepo.Create(ctx, link)

	// Test: Get existing link
	found, err := linkRepo.GetByID(ctx, link.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if found == nil {
		t.Fatal("GetByID() returned nil for existing link")
	}
	if found.Title != link.Title {
		t.Errorf("Title = %v, want %v", found.Title, link.Title)
	}

	// Test: Get non-existing link
	notFound, err := linkRepo.GetByID(ctx, 99999)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if notFound != nil {
		t.Error("GetByID() should return nil for non-existing link")
	}
}

func TestLinkRepository_GetActiveByUserID(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	user := createTestUser(t, userRepo, "activetest")

	// Create 2 active and 1 inactive link
	active1 := &model.Link{UserID: user.ID, Title: "Active 1", URL: "https://a.com", IsActive: true}
	active2 := &model.Link{UserID: user.ID, Title: "Active 2", URL: "https://b.com", IsActive: true}
	inactive := &model.Link{UserID: user.ID, Title: "Inactive", URL: "https://c.com", IsActive: false}

	linkRepo.Create(ctx, active1)
	linkRepo.Create(ctx, active2)
	linkRepo.Create(ctx, inactive)

	// Test: Get active links only
	activeLinks, err := linkRepo.GetActiveByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetActiveByUserID() error = %v", err)
	}
	if len(activeLinks) != 2 {
		t.Errorf("GetActiveByUserID() returned %d links, want 2", len(activeLinks))
	}
	for _, link := range activeLinks {
		if !link.IsActive {
			t.Errorf("GetActiveByUserID() returned inactive link: %s", link.Title)
		}
	}
}

func TestLinkRepository_Delete(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	user := createTestUser(t, userRepo, "deltest")

	link := &model.Link{UserID: user.ID, Title: "To Delete", URL: "https://del.com", IsActive: true}
	linkRepo.Create(ctx, link)

	err := linkRepo.Delete(ctx, link.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	found, _ := linkRepo.GetByID(ctx, link.ID)
	if found != nil {
		t.Error("Delete() did not remove link")
	}
}

func TestLinkRepository_UpdatePositions(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	user := createTestUser(t, userRepo, "reordertest")

	link1 := &model.Link{UserID: user.ID, Title: "Link 1", URL: "https://1.com", IsActive: true}
	link2 := &model.Link{UserID: user.ID, Title: "Link 2", URL: "https://2.com", IsActive: true}
	link3 := &model.Link{UserID: user.ID, Title: "Link 3", URL: "https://3.com", IsActive: true}

	linkRepo.Create(ctx, link1)
	linkRepo.Create(ctx, link2)
	linkRepo.Create(ctx, link3)

	// Reorder: Move link3 to first position
	positions := map[int64]int{
		link3.ID: 1,
		link1.ID: 2,
		link2.ID: 3,
	}

	err := linkRepo.UpdatePositions(ctx, user.ID, positions)
	if err != nil {
		t.Fatalf("UpdatePositions() error = %v", err)
	}

	// Verify new order
	links, _ := linkRepo.GetByUserID(ctx, user.ID)
	if links[0].ID != link3.ID {
		t.Errorf("First link should be link3, got ID %d", links[0].ID)
	}
	if links[1].ID != link1.ID {
		t.Errorf("Second link should be link1, got ID %d", links[1].ID)
	}
}

func TestLinkRepository_UserIsolation(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	user1 := createTestUser(t, userRepo, "user1")
	user2 := createTestUser(t, userRepo, "user2")

	link1 := &model.Link{UserID: user1.ID, Title: "User1 Link", URL: "https://u1.com", IsActive: true}
	link2 := &model.Link{UserID: user2.ID, Title: "User2 Link", URL: "https://u2.com", IsActive: true}

	linkRepo.Create(ctx, link1)
	linkRepo.Create(ctx, link2)

	// Each user should only see their own links
	user1Links, _ := linkRepo.GetByUserID(ctx, user1.ID)
	user2Links, _ := linkRepo.GetByUserID(ctx, user2.ID)

	if len(user1Links) != 1 {
		t.Errorf("User1 should have 1 link, got %d", len(user1Links))
	}
	if len(user2Links) != 1 {
		t.Errorf("User2 should have 1 link, got %d", len(user2Links))
	}
}
```

> 🧠 📱 "createTestUser helper — test-ൽ user create ചെയ്യുന്ന boilerplate reduce ചെയ്യാൻ. Every link test-ന് user വേണം."
> 📱 "PositionAutoIncrement — COALESCE(MAX(position), 0) + 1 logic work ചെയ്യുന്നുണ്ട് എന്ന് verify."
> 📱 "GetActiveByUserID — is_active filter. Public profile-ൽ inactive links കാണരുത്."
> 📱 "UpdatePositions — drag-drop reorder. Transaction-based update. Order correctly change ആയോ verify."
> 📱 "UserIsolation — ⚠️ important! User A-ന്റെ links User B-ന് access ചെയ്യാൻ പറ്റരുത്. Multi-tenant security test."

---

## Analytics Repository Tests

**⌨️ Create `internal/repository/analytics_test.go`:**
```go
package repository

import (
	"context"
	"testing"

	"linkbio/internal/model"
	"linkbio/internal/testutil"
)

func TestAnalyticsRepository_RecordPageView(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	analyticsRepo := NewAnalyticsRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "analyticstest", Email: "analytics@test.com",
		PasswordHash: "hash", DisplayName: "Analytics Test", Theme: "light",
	}
	userRepo.Create(ctx, user)

	err := analyticsRepo.RecordPageView(ctx, user.ID, "https://google.com", "Mozilla/5.0")
	if err != nil {
		t.Fatalf("RecordPageView() error = %v", err)
	}

	summary, err := analyticsRepo.GetSummary(ctx, user.ID, 1)
	if err != nil {
		t.Fatalf("GetSummary() error = %v", err)
	}
	if summary.TotalViews != 1 {
		t.Errorf("TotalViews = %d, want 1", summary.TotalViews)
	}
}

func TestAnalyticsRepository_RecordLinkClick(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	analyticsRepo := NewAnalyticsRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "clicktest", Email: "click@test.com",
		PasswordHash: "hash", DisplayName: "Click Test", Theme: "light",
	}
	userRepo.Create(ctx, user)

	link := &model.Link{UserID: user.ID, Title: "Test Link", URL: "https://example.com", IsActive: true}
	linkRepo.Create(ctx, link)

	err := analyticsRepo.RecordLinkClick(ctx, user.ID, link.ID, "https://instagram.com", "Mozilla/5.0")
	if err != nil {
		t.Fatalf("RecordLinkClick() error = %v", err)
	}

	summary, err := analyticsRepo.GetSummary(ctx, user.ID, 1)
	if err != nil {
		t.Fatalf("GetSummary() error = %v", err)
	}
	if summary.TotalClicks != 1 {
		t.Errorf("TotalClicks = %d, want 1", summary.TotalClicks)
	}
}

func TestAnalyticsRepository_GetSummary(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	analyticsRepo := NewAnalyticsRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "summarytest", Email: "summary@test.com",
		PasswordHash: "hash", DisplayName: "Summary Test", Theme: "light",
	}
	userRepo.Create(ctx, user)

	link := &model.Link{UserID: user.ID, Title: "Test Link", URL: "https://example.com", IsActive: true}
	linkRepo.Create(ctx, link)

	// Record multiple events
	for i := 0; i < 5; i++ {
		analyticsRepo.RecordPageView(ctx, user.ID, "", "")
	}
	for i := 0; i < 3; i++ {
		analyticsRepo.RecordLinkClick(ctx, user.ID, link.ID, "", "")
	}

	summary, err := analyticsRepo.GetSummary(ctx, user.ID, 7)
	if err != nil {
		t.Fatalf("GetSummary() error = %v", err)
	}
	if summary.TotalViews != 5 {
		t.Errorf("TotalViews = %d, want 5", summary.TotalViews)
	}
	if summary.TotalClicks != 3 {
		t.Errorf("TotalClicks = %d, want 3", summary.TotalClicks)
	}
}

func TestAnalyticsRepository_UserIsolation(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	analyticsRepo := NewAnalyticsRepository(db)
	ctx := context.Background()

	user1 := &model.User{
		Username: "iso1", Email: "iso1@test.com",
		PasswordHash: "hash", DisplayName: "Iso 1", Theme: "light",
	}
	user2 := &model.User{
		Username: "iso2", Email: "iso2@test.com",
		PasswordHash: "hash", DisplayName: "Iso 2", Theme: "light",
	}
	userRepo.Create(ctx, user1)
	userRepo.Create(ctx, user2)

	for i := 0; i < 10; i++ {
		analyticsRepo.RecordPageView(ctx, user1.ID, "", "")
	}
	for i := 0; i < 3; i++ {
		analyticsRepo.RecordPageView(ctx, user2.ID, "", "")
	}

	summary1, _ := analyticsRepo.GetSummary(ctx, user1.ID, 7)
	summary2, _ := analyticsRepo.GetSummary(ctx, user2.ID, 7)

	if summary1.TotalViews != 10 {
		t.Errorf("User1 TotalViews = %d, want 10", summary1.TotalViews)
	}
	if summary2.TotalViews != 3 {
		t.Errorf("User2 TotalViews = %d, want 3", summary2.TotalViews)
	}
}
```

> 🧠 📱 "RecordPageView, RecordLinkClick — individual function tests. ഓരോ analytics function correct ആണോ verify."
> 📱 "GetSummary — multiple events record ചെയ്ത ശേഷം aggregation correct ആണോ check."
> 📱 "UserIsolation — same pattern as link tests. User A-ന്റെ analytics User B-ന് visible ആകരുത്."

---

## Health Handler Test

**⌨️ Create `internal/handler/health_test.go`:**
```go
package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler_Check(t *testing.T) {
	h := &HealthHandler{}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	h.Check(rec, req)

	// Check status code
	if rec.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
	}

	// Check content type
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %s, want application/json", contentType)
	}

	// Check response body
	var response map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if response["status"] != "ok" {
		t.Errorf("status = %s, want ok", response["status"])
	}
}

func TestHealthHandler_Check_Methods(t *testing.T) {
	h := &HealthHandler{}

	methods := []string{
		http.MethodGet, http.MethodPost,
		http.MethodPut, http.MethodDelete,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/health", nil)
			rec := httptest.NewRecorder()
			h.Check(rec, req)
			if rec.Code != http.StatusOK {
				t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
			}
		})
	}
}

func BenchmarkHealthHandler_Check(b *testing.B) {
	h := &HealthHandler{}
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		h.Check(rec, req)
	}
}
```

> 🧠 📱 "httptest.NewRequest — fake HTTP request. httptest.NewRecorder — response capture."
> 📱 "Real server start ചെയ്യാതെ handler test ചെയ്യാം!"
> 📱 "Check_Methods — HealthHandler ഏത് HTTP method-നും work ചെയ്യണം. GET, POST, PUT, DELETE — everything."

---

## Auth Handler Tests

> 📱 "ഇനി auth handler tests — login, register, logout. ഇത് comprehensive ആണ്. ⚠️ bcrypt use ചെയ്യുന്നതുകൊണ്ട് slightly slow, but thorough."

**⌨️ Create `internal/handler/auth_test.go`:**
```go
package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"linkbio/internal/model"
	"linkbio/internal/pkg/response"
	"linkbio/internal/repository"
	"linkbio/internal/testutil"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func setupAuthHandler(t *testing.T) (*AuthHandler, *repository.UserRepository) {
	t.Helper()

	db := testutil.TestDB(t)
	log := testutil.TestLogger()

	userRepo := repository.NewUserRepository(db)
	resp := response.New(log)
	store := sessions.NewCookieStore([]byte("test-secret-key-32-chars-minimum!"))

	h := &AuthHandler{
		log:      log,
		resp:     resp,
		store:    store,
		userRepo: userRepo,
	}

	return h, userRepo
}

func TestAuthHandler_Login_Success(t *testing.T) {
	h, userRepo := setupAuthHandler(t)
	ctx := context.Background()

	password := "testpass123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &model.User{
		Username: "loginuser", Email: "login@test.com",
		PasswordHash: string(hash), DisplayName: "Login User", Theme: "light",
	}
	userRepo.Create(ctx, user)

	form := url.Values{}
	form.Add("email", "login@test.com")
	form.Add("password", password)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Header().Get("HX-Redirect") != "/dashboard" {
		t.Errorf("HX-Redirect = %s, want /dashboard", rec.Header().Get("HX-Redirect"))
	}
}

func TestAuthHandler_Login_WrongPassword(t *testing.T) {
	h, userRepo := setupAuthHandler(t)
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
	user := &model.User{
		Username: "wrongpass", Email: "wrong@test.com",
		PasswordHash: string(hash), DisplayName: "Wrong Pass", Theme: "light",
	}
	userRepo.Create(ctx, user)

	form := url.Values{}
	form.Add("email", "wrong@test.com")
	form.Add("password", "wrongpassword")

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthHandler_Login_UserNotFound(t *testing.T) {
	h, _ := setupAuthHandler(t)

	form := url.Values{}
	form.Add("email", "notexists@test.com")
	form.Add("password", "anypassword")

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthHandler_Login_MissingFields(t *testing.T) {
	h, _ := setupAuthHandler(t)

	tests := []struct {
		name     string
		email    string
		password string
	}{
		{"missing email", "", "password"},
		{"missing password", "test@test.com", ""},
		{"missing both", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("email", tt.email)
			form.Add("password", tt.password)

			req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()

			h.Login(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("Status = %d, want %d", rec.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestAuthHandler_Register_Success(t *testing.T) {
	h, userRepo := setupAuthHandler(t)
	ctx := context.Background()

	form := url.Values{}
	form.Add("username", "newuser")
	form.Add("email", "new@test.com")
	form.Add("password", "password123")

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Register(rec, req)

	if rec.Header().Get("HX-Redirect") != "/dashboard" {
		t.Errorf("HX-Redirect = %s, want /dashboard", rec.Header().Get("HX-Redirect"))
	}

	// Verify user was created
	user, _ := userRepo.GetByUsername(ctx, "newuser")
	if user == nil {
		t.Error("User was not created")
	}
	// Verify password was hashed (not stored plain)
	if user.PasswordHash == "password123" {
		t.Error("Password was not hashed")
	}
}

func TestAuthHandler_Register_DuplicateUsername(t *testing.T) {
	h, userRepo := setupAuthHandler(t)
	ctx := context.Background()

	existing := &model.User{
		Username: "existing", Email: "existing@test.com",
		PasswordHash: "hash", DisplayName: "Existing", Theme: "light",
	}
	userRepo.Create(ctx, existing)

	form := url.Values{}
	form.Add("username", "existing")
	form.Add("email", "different@test.com")
	form.Add("password", "password123")

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Register(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusConflict)
	}
}

func TestAuthHandler_Register_ShortPassword(t *testing.T) {
	h, _ := setupAuthHandler(t)

	form := url.Values{}
	form.Add("username", "shortpass")
	form.Add("email", "short@test.com")
	form.Add("password", "12345") // Less than 6 chars

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	h, _ := setupAuthHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	rec := httptest.NewRecorder()

	h.Logout(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusSeeOther)
	}
	if rec.Header().Get("Location") != "/" {
		t.Errorf("Location = %s, want /", rec.Header().Get("Location"))
	}
}
```

> 🧠 📱 "setupAuthHandler — test setup helper. In-memory DB, test logger, gorilla session store with test secret. Handler struct directly create ചെയ്യുന്നു."
> 📱 "Login_Success — bcrypt hash create, form POST, HX-Redirect header check. Full HTMX flow test."
> 📱 "Login_WrongPassword — bcrypt compare fail → 401 Unauthorized."
> 📱 "Login_MissingFields — table-driven. Empty email, empty password, both empty — all 400 BadRequest."
> 📱 "Register_Success — user actually created ആയോ verify. Password hash ആയോ (plain text അല്ല) verify."
> 📱 "Register_DuplicateUsername — 409 Conflict. Existing username-ൽ register ചെയ്യാൻ attempt."
> 📱 "Register_ShortPassword — 6 characters minimum. '12345' → 400 BadRequest."
> 📱 "Logout — session MaxAge = -1, redirect to '/'. Cookie delete ചെയ്യുന്നു."

---

## Benchmarks

> 📱 "ഇനി benchmarks — performance measure ചെയ്യാൻ. ⚠️ benchmarks-ൽ `*testing.B` use ചെയ്യുന്നു, `*testing.T` അല്ല."

**⌨️ Create `internal/repository/benchmark_test.go`:**
```go
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"linkbio/internal/model"

	_ "modernc.org/sqlite"
)

// benchDB creates a fresh in-memory database for benchmarks
func benchDB(b *testing.B) *sql.DB {
	b.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		b.Fatalf("failed to open bench db: %v", err)
	}

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	if err := Migrate(db, log); err != nil {
		b.Fatalf("failed to migrate bench db: %v", err)
	}

	return db
}

func BenchmarkUserRepository_Create(b *testing.B) {
	db := benchDB(b)
	defer db.Close()
	repo := NewUserRepository(db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &model.User{
			Username:     fmt.Sprintf("user%d", i),
			Email:        fmt.Sprintf("user%d@test.com", i),
			PasswordHash: "hash",
			DisplayName:  "Bench User",
			Theme:        "light",
		}
		repo.Create(ctx, user)
	}
}

func BenchmarkUserRepository_GetByEmail(b *testing.B) {
	db := benchDB(b)
	defer db.Close()
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "findme", Email: "findme@test.com",
		PasswordHash: "hash", DisplayName: "Find Me", Theme: "light",
	}
	repo.Create(ctx, user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.GetByEmail(ctx, "findme@test.com")
	}
}

func BenchmarkLinkRepository_GetByUserID(b *testing.B) {
	db := benchDB(b)
	defer db.Close()
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "getlinks", Email: "getlinks@test.com",
		PasswordHash: "hash", DisplayName: "Get Links", Theme: "light",
	}
	userRepo.Create(ctx, user)

	for i := 0; i < 10; i++ {
		link := &model.Link{
			UserID: user.ID, Title: fmt.Sprintf("Link %d", i),
			URL: "https://example.com", IsActive: true,
		}
		linkRepo.Create(ctx, link)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		linkRepo.GetByUserID(ctx, user.ID)
	}
}

func BenchmarkAnalyticsRepository_RecordPageView(b *testing.B) {
	db := benchDB(b)
	defer db.Close()
	userRepo := NewUserRepository(db)
	analyticsRepo := NewAnalyticsRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "analyticsbench", Email: "analyticsbench@test.com",
		PasswordHash: "hash", DisplayName: "Analytics Bench", Theme: "light",
	}
	userRepo.Create(ctx, user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyticsRepo.RecordPageView(ctx, user.ID, "https://google.com", "Mozilla/5.0")
	}
}
```

> 🧠 📱 "benchDB helper — benchmarks-ൽ `*testing.T` pass ചെയ്യാൻ പറ്റില്ല. `*testing.B` accept ചെയ്യുന്ന separate helper."
> 📱 "Migrate(db, log) — db.go-ൽ ഉള്ള real migration function use ചെയ്യുന്നു. testutil-ൽ inline migrations-ൽ നിന്ന് different approach — benchmarks same package ആയതുകൊണ്ട് direct access."
> 📱 "b.ResetTimer() — setup time exclude. DB create, user insert — ഇതൊന്നും benchmark-ൽ count ആകില്ല."
> 📱 "b.N — Go automatically decides iterations. Fast function = more iterations. Accurate timing."

---

## Running Tests

```bash
# All tests
go test ./...

# Verbose output — see each test name
go test -v ./...

# Specific test function
go test -v -run TestUserRepository ./internal/repository

# Coverage percentage
go test -cover ./...

# HTML coverage report
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

# Benchmarks with memory stats
go test -bench=. -benchmem ./...

# Race condition detection
go test -race ./...
```

> 🧠 📱 "Expected benchmark output:"
```
BenchmarkUserRepository_GetByEmail-8     100000    15234 ns/op    512 B/op    8 allocs/op
```
> 📱 "100000 = iterations. 15234 ns = ~15 microseconds per query. 512 B = memory per op. 8 allocs = heap allocations."

> 📱 "`make test` — Makefile shortcut. `go test -v ./...` run ചെയ്യും."

---

## Test Summary

| Package | Tests | Benchmarks | What's Tested |
|---------|-------|------------|--------------|
| `repository` | 17 | 4+ | User CRUD, Link CRUD, Analytics, position auto-increment, user isolation |
| `handler` | 13 | 1 | Health check, Login (success/fail/missing), Register (success/duplicate/short password), Logout |
| **Total** | **30** | **5+** | |

---

## Coverage Goals

| Package | Target |
|---------|--------|
| repository | 80%+ |
| handler | 70%+ |
| middleware | 60%+ |

> 📱 "100% coverage വേണ്ട. Critical paths cover ചെയ്താൽ മതി."

---

## Testing Best Practices

1. **AAA Pattern**: Arrange → Act → Assert
2. **Table-Driven**: Multiple cases, one function — `tests := []struct{}`
3. **Isolation**: Each test independent — own `testutil.TestDB(t)`
4. **Names**: `TestUserRepository_Create`, `TestAuthHandler_Login_WrongPassword`
5. **Edge Cases**: Empty input, nil, duplicate data, wrong passwords
6. **Helpers**: `createTestUser()`, `setupAuthHandler()` — reduce boilerplate
7. **User Isolation**: Always test that User A can't see User B's data
