package repository

import (
	"context"
	"testing"

	"linkbio/internal/model"
	"linkbio/internal/testutil"
)

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
		UserID:   user.ID,
		Title:    "My Website",
		URL:      "https://example.com",
		IsActive: true,
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

	// Create 3 links
	for i := 1; i <= 3; i++ {
		link := &model.Link{
			UserID:   user.ID,
			Title:    "Link",
			URL:      "https://example.com",
			IsActive: true,
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
		UserID:   user.ID,
		Title:    "Test Link",
		URL:      "https://test.com",
		Icon:     "🔗",
		IsActive: true,
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
	if found.IsActive != true {
		t.Error("IsActive should be true")
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

func TestLinkRepository_GetByUserID(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	user := createTestUser(t, userRepo, "getbyuseridtest")

	// Create 3 links
	for i := 1; i <= 3; i++ {
		link := &model.Link{
			UserID:   user.ID,
			Title:    "Link " + string(rune('0'+i)),
			URL:      "https://example.com",
			IsActive: true,
		}
		linkRepo.Create(ctx, link)
	}

	// Test: Get all links
	links, err := linkRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByUserID() error = %v", err)
	}
	if len(links) != 3 {
		t.Errorf("GetByUserID() returned %d links, want 3", len(links))
	}

	// Verify order by position
	for i, link := range links {
		if link.Position != i+1 {
			t.Errorf("Link at index %d has position %d, want %d", i, link.Position, i+1)
		}
	}
}

func TestLinkRepository_GetActiveByUserID(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	user := createTestUser(t, userRepo, "activetest")

	// Create 2 active and 1 inactive link
	activeLink1 := &model.Link{UserID: user.ID, Title: "Active 1", URL: "https://a.com", IsActive: true}
	activeLink2 := &model.Link{UserID: user.ID, Title: "Active 2", URL: "https://b.com", IsActive: true}
	inactiveLink := &model.Link{UserID: user.ID, Title: "Inactive", URL: "https://c.com", IsActive: false}

	linkRepo.Create(ctx, activeLink1)
	linkRepo.Create(ctx, activeLink2)
	linkRepo.Create(ctx, inactiveLink)

	// Test: Get active links only
	activeLinks, err := linkRepo.GetActiveByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetActiveByUserID() error = %v", err)
	}

	if len(activeLinks) != 2 {
		t.Errorf("GetActiveByUserID() returned %d links, want 2", len(activeLinks))
	}

	// Verify all returned links are active
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

	// Test: Delete link
	err := linkRepo.Delete(ctx, link.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify: Link no longer exists
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

	// Create 3 links
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
	if links[2].ID != link2.ID {
		t.Errorf("Third link should be link2, got ID %d", links[2].ID)
	}
}

func TestLinkRepository_UserIsolation(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	// Create two users
	user1 := createTestUser(t, userRepo, "user1")
	user2 := createTestUser(t, userRepo, "user2")

	// Create links for each user
	link1 := &model.Link{UserID: user1.ID, Title: "User1 Link", URL: "https://u1.com", IsActive: true}
	link2 := &model.Link{UserID: user2.ID, Title: "User2 Link", URL: "https://u2.com", IsActive: true}

	linkRepo.Create(ctx, link1)
	linkRepo.Create(ctx, link2)

	// Test: Each user should only see their own links
	user1Links, _ := linkRepo.GetByUserID(ctx, user1.ID)
	user2Links, _ := linkRepo.GetByUserID(ctx, user2.ID)

	if len(user1Links) != 1 {
		t.Errorf("User1 should have 1 link, got %d", len(user1Links))
	}
	if len(user2Links) != 1 {
		t.Errorf("User2 should have 1 link, got %d", len(user2Links))
	}
	if user1Links[0].Title != "User1 Link" {
		t.Error("User1 got wrong link")
	}
	if user2Links[0].Title != "User2 Link" {
		t.Error("User2 got wrong link")
	}
}
