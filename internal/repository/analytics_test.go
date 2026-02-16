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

	// Setup: Create user
	user := &model.User{
		Username:     "analyticstest",
		Email:        "analytics@test.com",
		PasswordHash: "hash",
		DisplayName:  "Analytics Test",
		Theme:        "light",
	}
	userRepo.Create(ctx, user)

	// Test: Record page view
	err := analyticsRepo.RecordPageView(ctx, user.ID, "https://google.com", "Mozilla/5.0")
	if err != nil {
		t.Fatalf("RecordPageView() error = %v", err)
	}

	// Verify via summary
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

	// Setup: Create user and link
	user := &model.User{
		Username:     "clicktest",
		Email:        "click@test.com",
		PasswordHash: "hash",
		DisplayName:  "Click Test",
		Theme:        "light",
	}
	userRepo.Create(ctx, user)

	link := &model.Link{
		UserID:   user.ID,
		Title:    "Test Link",
		URL:      "https://example.com",
		IsActive: true,
	}
	linkRepo.Create(ctx, link)

	// Test: Record link click
	err := analyticsRepo.RecordLinkClick(ctx, user.ID, link.ID, "https://instagram.com", "Mozilla/5.0")
	if err != nil {
		t.Fatalf("RecordLinkClick() error = %v", err)
	}

	// Verify via summary
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

	// Setup: Create user and link
	user := &model.User{
		Username:     "summarytest",
		Email:        "summary@test.com",
		PasswordHash: "hash",
		DisplayName:  "Summary Test",
		Theme:        "light",
	}
	userRepo.Create(ctx, user)

	link := &model.Link{
		UserID:   user.ID,
		Title:    "Test Link",
		URL:      "https://example.com",
		IsActive: true,
	}
	linkRepo.Create(ctx, link)

	// Record multiple events
	for i := 0; i < 5; i++ {
		analyticsRepo.RecordPageView(ctx, user.ID, "", "")
	}
	for i := 0; i < 3; i++ {
		analyticsRepo.RecordLinkClick(ctx, user.ID, link.ID, "", "")
	}

	// Test: Get summary
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

	// Create two users
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

	// Record different amounts for each user
	for i := 0; i < 10; i++ {
		analyticsRepo.RecordPageView(ctx, user1.ID, "", "")
	}
	for i := 0; i < 3; i++ {
		analyticsRepo.RecordPageView(ctx, user2.ID, "", "")
	}

	// Test: Each user sees only their own analytics
	summary1, _ := analyticsRepo.GetSummary(ctx, user1.ID, 7)
	summary2, _ := analyticsRepo.GetSummary(ctx, user2.ID, 7)

	if summary1.TotalViews != 10 {
		t.Errorf("User1 TotalViews = %d, want 10", summary1.TotalViews)
	}
	if summary2.TotalViews != 3 {
		t.Errorf("User2 TotalViews = %d, want 3", summary2.TotalViews)
	}
}
