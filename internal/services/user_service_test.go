package services_test

import (
	"context"
	"testing"

	"github.com/thbappy7706/go-inertia-starter-kit/internal/repositories"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/services"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/testhelper"
)

func TestUserService_List_Search_Pagination(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo)
	ctx := context.Background()

	// Seed 3 users
	_, _, _ = authService.Register(ctx, "Alice Smith", "alice@example.com", "password123", "password123")
	_, _, _ = authService.Register(ctx, "Bob Jones", "bob@example.com", "password123", "password123")
	_, _, _ = authService.Register(ctx, "Charlie Brown", "charlie@example.com", "password123", "password123")

	// 1. Listing all users
	allUsers, err := userService.List(ctx, "", 1, 10)
	if err != nil {
		t.Fatalf("list users error: %v", err)
	}
	if allUsers.Total != 3 {
		t.Errorf("expected 3 users, got %d", allUsers.Total)
	}

	// 2. Search
	searchResult, err := userService.List(ctx, "Alice", 1, 10)
	if err != nil {
		t.Fatalf("search users error: %v", err)
	}
	if searchResult.Total != 1 || searchResult.Data[0].Name != "Alice Smith" {
		t.Errorf("search result mismatch: %v", searchResult)
	}

	// 3. Pagination
	pagedResult, err := userService.List(ctx, "", 1, 2)
	if err != nil {
		t.Fatalf("pagination error: %v", err)
	}
	if len(pagedResult.Data) != 2 || pagedResult.LastPage != 2 {
		t.Errorf("pagination mismatch: len %d, lastPage %d", len(pagedResult.Data), pagedResult.LastPage)
	}
}
