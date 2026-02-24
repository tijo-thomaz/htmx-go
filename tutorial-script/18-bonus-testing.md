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
func TestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// IMPORTANT: SQLite :memory: creates separate DB per connection
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Enable foreign keys
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

	t.Cleanup(func() { db.Close() })
	return db
}

func TestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}
```

> 🧠 📱 "t.Helper() — test fail line number correct ആക്കും."
> 📱 ":memory: — in-memory SQLite. Test-ന് ശേഷം auto-delete."
> 📱 "t.Cleanup() — test end-ൽ db.Close() auto-call."
> 📱 "⚠️ Import cycle avoid ചെയ്യാൻ migrations inline ആണ്. testutil package-ൽ repository import ചെയ്താൽ cycle ആകും!"
> 📱 "SetMaxOpenConns(1) — SQLite :memory: database ഓരോ connection-നും separate ആണ്. One connection = consistent data."
> 📱 "PRAGMA foreign_keys = ON — SQLite default-ൽ foreign keys OFF ആണ്! Tests-ൽ ON ചെയ്യണം."

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
		Username: "testuser", Email: "test@example.com",
		PasswordHash: "hashedpassword", DisplayName: "Test User", Theme: "light",
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if user.ID == 0 {
		t.Error("Create() did not set user ID")
	}

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

	user := &model.User{
		Username: "emailtest", Email: "find@example.com",
		PasswordHash: "hash", DisplayName: "Email Test", Theme: "dark",
	}
	repo.Create(ctx, user)

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
			if (found != nil) != tt.wantFound {
				t.Errorf("GetByEmail() found = %v, want %v", found != nil, tt.wantFound)
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
	repo.Create(ctx, user1)

	user2 := &model.User{
		Username: "duplicate", Email: "user2@example.com",
		PasswordHash: "hash", DisplayName: "User 2", Theme: "light",
	}
	err := repo.Create(ctx, user2)
	if err == nil {
		t.Error("Create() should fail for duplicate username")
	}
}
```

> 🧠 📱 "Table-driven tests — tests := []struct{}. Multiple cases, one function. Go community standard."
> 📱 "t.Run() — sub-tests. ഓരോ case independently run."

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

func TestLinkRepository_Create(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "linktest", Email: "link@test.com",
		PasswordHash: "hash", DisplayName: "Link Test", Theme: "light",
	}
	userRepo.Create(ctx, user)

	link := &model.Link{UserID: user.ID, Title: "My Website", URL: "https://example.com", IsActive: true}
	err := linkRepo.Create(ctx, link)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if link.ID == 0 {
		t.Error("Create() did not set link ID")
	}
	if link.Position != 1 {
		t.Errorf("Position = %v, want 1", link.Position)
	}
}

func TestLinkRepository_Delete(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "deltest", Email: "del@test.com",
		PasswordHash: "hash", DisplayName: "Del Test", Theme: "light",
	}
	userRepo.Create(ctx, user)

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
```

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

	if rec.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusOK)
	}

	var response map[string]string
	json.Unmarshal(rec.Body.Bytes(), &response)
	if response["status"] != "ok" {
		t.Errorf("status = %s, want ok", response["status"])
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

---

## Benchmarks

**⌨️ Create `internal/repository/benchmark_test.go`:**
```go
package repository

import (
	"context"
	"fmt"
	"testing"
	"linkbio/internal/model"
	"linkbio/internal/testutil"
)

func BenchmarkUserRepository_GetByEmail(b *testing.B) {
	db := testutil.TestDB(&testing.T{})
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
	db := testutil.TestDB(&testing.T{})
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "linkbench", Email: "linkbench@test.com",
		PasswordHash: "hash", DisplayName: "Link Bench", Theme: "light",
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
```

> 🧠 📱 "b.N — Go automatically decides iterations."
> 📱 "b.ResetTimer() — setup time exclude."

---

## Running Tests

```bash
go test ./...                          # All tests
go test -v ./...                       # Verbose
go test -v -run TestUserRepository ./internal/repository  # Specific
go test -cover ./...                   # Coverage %
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out  # HTML report
go test -bench=. -benchmem ./...       # Benchmarks
go test -race ./...                    # Race condition detection
```

> 🧠 📱 "Expected benchmark output:"
```
BenchmarkUserRepository_GetByEmail-8     100000    15234 ns/op    512 B/op    8 allocs/op
```
> 📱 "100000 = iterations. 15234 ns = ~15 microseconds per query. 512 B = memory per op."

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
2. **Table-Driven**: Multiple cases, one function
3. **Isolation**: Each test independent
4. **Names**: `TestUserRepository_Create_DuplicateEmail`
5. **Edge Cases**: Empty input, nil, boundaries
