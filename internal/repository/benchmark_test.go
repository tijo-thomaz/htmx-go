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

	// Setup: Create user to find
	user := &model.User{
		Username:     "findme",
		Email:        "findme@test.com",
		PasswordHash: "hash",
		DisplayName:  "Find Me",
		Theme:        "light",
	}
	repo.Create(ctx, user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.GetByEmail(ctx, "findme@test.com")
	}
}

func BenchmarkUserRepository_GetByUsername(b *testing.B) {
	db := benchDB(b)
	defer db.Close()
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Setup: Create user to find
	user := &model.User{
		Username:     "benchuser",
		Email:        "benchuser@test.com",
		PasswordHash: "hash",
		DisplayName:  "Bench User",
		Theme:        "light",
	}
	repo.Create(ctx, user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.GetByUsername(ctx, "benchuser")
	}
}

func BenchmarkLinkRepository_Create(b *testing.B) {
	db := benchDB(b)
	defer db.Close()
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	// Setup: Create user
	user := &model.User{
		Username: "linkbench", Email: "linkbench@test.com",
		PasswordHash: "hash", DisplayName: "Link Bench", Theme: "light",
	}
	userRepo.Create(ctx, user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		link := &model.Link{
			UserID:   user.ID,
			Title:    "Bench Link",
			URL:      "https://example.com",
			IsActive: true,
		}
		linkRepo.Create(ctx, link)
	}
}

func BenchmarkLinkRepository_GetByUserID(b *testing.B) {
	db := benchDB(b)
	defer db.Close()
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	// Setup: Create user with 10 links
	user := &model.User{
		Username: "getlinks", Email: "getlinks@test.com",
		PasswordHash: "hash", DisplayName: "Get Links", Theme: "light",
	}
	userRepo.Create(ctx, user)

	for i := 0; i < 10; i++ {
		link := &model.Link{
			UserID:   user.ID,
			Title:    fmt.Sprintf("Link %d", i),
			URL:      "https://example.com",
			IsActive: true,
		}
		linkRepo.Create(ctx, link)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		linkRepo.GetByUserID(ctx, user.ID)
	}
}

func BenchmarkLinkRepository_GetActiveByUserID(b *testing.B) {
	db := benchDB(b)
	defer db.Close()
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	// Setup: Create user with 10 links (8 active, 2 inactive)
	user := &model.User{
		Username: "activelinks", Email: "activelinks@test.com",
		PasswordHash: "hash", DisplayName: "Active Links", Theme: "light",
	}
	userRepo.Create(ctx, user)

	for i := 0; i < 10; i++ {
		link := &model.Link{
			UserID:   user.ID,
			Title:    fmt.Sprintf("Link %d", i),
			URL:      "https://example.com",
			IsActive: i < 8, // First 8 are active
		}
		linkRepo.Create(ctx, link)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		linkRepo.GetActiveByUserID(ctx, user.ID)
	}
}

func BenchmarkAnalyticsRepository_RecordPageView(b *testing.B) {
	db := benchDB(b)
	defer db.Close()
	userRepo := NewUserRepository(db)
	analyticsRepo := NewAnalyticsRepository(db)
	ctx := context.Background()

	// Setup: Create user
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

func BenchmarkAnalyticsRepository_GetSummary(b *testing.B) {
	db := benchDB(b)
	defer db.Close()
	userRepo := NewUserRepository(db)
	analyticsRepo := NewAnalyticsRepository(db)
	ctx := context.Background()

	// Setup: Create user with some analytics data
	user := &model.User{
		Username: "summarybench", Email: "summarybench@test.com",
		PasswordHash: "hash", DisplayName: "Summary Bench", Theme: "light",
	}
	userRepo.Create(ctx, user)

	// Add 100 page views
	for i := 0; i < 100; i++ {
		analyticsRepo.RecordPageView(ctx, user.ID, "", "")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyticsRepo.GetSummary(ctx, user.ID, 28)
	}
}
