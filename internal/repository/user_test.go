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
		{
			name:      "existing email",
			email:     "find@example.com",
			wantFound: true,
		},
		{
			name:      "non-existing email",
			email:     "notfound@example.com",
			wantFound: false,
		},
		{
			name:      "empty email",
			email:     "",
			wantFound: false,
		},
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

	// Setup: Create a user
	user := &model.User{
		Username:     "findbyname",
		Email:        "findbyname@example.com",
		PasswordHash: "hash",
		DisplayName:  "Find By Name",
		Theme:        "light",
	}
	repo.Create(ctx, user)

	// Test cases
	tests := []struct {
		name      string
		username  string
		wantFound bool
	}{
		{
			name:      "existing username",
			username:  "findbyname",
			wantFound: true,
		},
		{
			name:      "non-existing username",
			username:  "notfound",
			wantFound: false,
		},
		{
			name:      "case sensitive",
			username:  "FINDBYNAME",
			wantFound: false, // SQLite is case-sensitive by default
		},
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

	// Create first user
	user1 := &model.User{
		Username:     "duplicate",
		Email:        "user1@example.com",
		PasswordHash: "hash",
		DisplayName:  "User 1",
		Theme:        "light",
	}
	if err := repo.Create(ctx, user1); err != nil {
		t.Fatalf("Create() first user error = %v", err)
	}

	// Try to create second user with same username
	user2 := &model.User{
		Username:     "duplicate", // Same username!
		Email:        "user2@example.com",
		PasswordHash: "hash",
		DisplayName:  "User 2",
		Theme:        "light",
	}
	err := repo.Create(ctx, user2)

	// Should fail with unique constraint error
	if err == nil {
		t.Error("Create() should fail for duplicate username")
	}
}

func TestUserRepository_DuplicateEmail(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create first user
	user1 := &model.User{
		Username:     "user1",
		Email:        "duplicate@example.com",
		PasswordHash: "hash",
		DisplayName:  "User 1",
		Theme:        "light",
	}
	if err := repo.Create(ctx, user1); err != nil {
		t.Fatalf("Create() first user error = %v", err)
	}

	// Try to create second user with same email
	user2 := &model.User{
		Username:     "user2",
		Email:        "duplicate@example.com", // Same email!
		PasswordHash: "hash",
		DisplayName:  "User 2",
		Theme:        "light",
	}
	err := repo.Create(ctx, user2)

	// Should fail with unique constraint error
	if err == nil {
		t.Error("Create() should fail for duplicate email")
	}
}
