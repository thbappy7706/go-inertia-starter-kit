package services_test

import (
	"context"
	"testing"

	"github.com/thbappy7706/go-inertia-starter-kit/internal/repositories"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/services"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/testhelper"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register_Success(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	ctx := context.Background()

	user, valErrors, err := authService.Register(ctx, "Test User", "test@example.com", "password123", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(valErrors) > 0 {
		t.Fatalf("unexpected validation errors: %v", valErrors)
	}
	if user == nil || user.ID == 0 {
		t.Fatalf("expected created user with ID, got %v", user)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", user.Email)
	}

	// Password must be securely hashed
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123")); err != nil {
		t.Errorf("password was not hashed properly: %v", err)
	}
}

func TestAuthService_Register_ValidationFailures(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	ctx := context.Background()

	// Short name
	_, valErrors, _ := authService.Register(ctx, "a", "valid@example.com", "password123", "password123")
	if valErrors["name"] == "" {
		t.Errorf("expected name validation error, got none")
	}

	// Invalid email
	_, valErrors, _ = authService.Register(ctx, "Valid Name", "invalid-email", "password123", "password123")
	if valErrors["email"] == "" {
		t.Errorf("expected email validation error, got none")
	}

	// Short password
	_, valErrors, _ = authService.Register(ctx, "Valid Name", "valid@example.com", "short", "short")
	if valErrors["password"] == "" {
		t.Errorf("expected password validation error, got none")
	}

	// Password mismatch
	_, valErrors, _ = authService.Register(ctx, "Valid Name", "valid@example.com", "password123", "different123")
	if valErrors["password_confirmation"] == "" {
		t.Errorf("expected password_confirmation validation error, got none")
	}
}

func TestAuthService_Register_DuplicateEmailRejected(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	ctx := context.Background()

	_, _, err := authService.Register(ctx, "User One", "dup@example.com", "password123", "password123")
	if err != nil {
		t.Fatalf("unexpected error creating first user: %v", err)
	}

	_, valErrors, err := authService.Register(ctx, "User Two", "dup@example.com", "password123", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if valErrors["email"] == "" {
		t.Errorf("expected duplicate email error, got none")
	}
}

func TestAuthService_Login(t *testing.T) {
	db, cleanup := testhelper.SetupTestDB(t)
	defer cleanup()

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	ctx := context.Background()

	_, _, err := authService.Register(ctx, "Login User", "login@example.com", "secret1234", "secret1234")
	if err != nil {
		t.Fatalf("unexpected error registering user: %v", err)
	}

	// Login success
	user, valErrors, err := authService.Login(ctx, "login@example.com", "secret1234")
	if err != nil {
		t.Fatalf("login error: %v", err)
	}
	if len(valErrors) > 0 {
		t.Fatalf("unexpected val errors on valid login: %v", valErrors)
	}
	if user == nil || user.Email != "login@example.com" {
		t.Errorf("expected user email login@example.com, got %v", user)
	}

	// Login invalid password
	_, valErrors, _ = authService.Login(ctx, "login@example.com", "wrongpass")
	if valErrors["email"] == "" {
		t.Errorf("expected generic error on wrong password, got none")
	}

	// Login non-existent email
	_, valErrors, _ = authService.Login(ctx, "nobody@example.com", "anypass")
	if valErrors["email"] == "" {
		t.Errorf("expected generic error on non-existent user, got none")
	}
}
