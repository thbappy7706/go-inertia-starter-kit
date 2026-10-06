package routes_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gonertia "github.com/romsar/gonertia/v3"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/config"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/handlers"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/middleware"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/repositories"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/services"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/session"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/testhelper"
	"github.com/thbappy7706/go-inertia-starter-kit/routes"
)

func setupTestRouter(t *testing.T) (*http.ServeMux, http.Handler, *session.Manager, repositories.UserRepository) {
	t.Helper()
	db, _ := testhelper.SetupTestDB(t)

	cfg := &config.Config{
		AppEnv:        "testing",
		SessionSecret: "test-secret-key-32-chars-long-1234",
	}

	sm := session.NewManager(cfg)
	flashProvider := session.NewFlashProvider(sm.SessionManager)

	i, err := gonertia.New("<div id=\"app\" data-page='{{ .page }}'></div>",
		gonertia.WithFlashProvider(flashProvider),
	)
	if err != nil {
		t.Fatalf("failed to init inertia: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	productRepo := repositories.NewProductRepository(db)

	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	productService := services.NewProductService(productRepo)

	authMid := middleware.NewAuthMiddleware(i, sm.SessionManager, userRepo)
	authHandler := handlers.NewAuthHandler(i, sm.SessionManager, authService)
	dashboardHandler := handlers.NewDashboardHandler(i, productService, userRepo)
	userHandler := handlers.NewUserHandler(i, userService)
	productHandler := handlers.NewProductHandler(i, sm.SessionManager, productService)

	router := routes.Setup(routes.RouterParams{
		Inertia:          i,
		SessionManager:   sm.SessionManager,
		AuthMiddleware:   authMid,
		AuthHandler:      authHandler,
		DashboardHandler: dashboardHandler,
		UserHandler:      userHandler,
		ProductHandler:   productHandler,
	})

	return nil, router, sm, userRepo
}

func TestRoute_ProtectedRequiresAuth(t *testing.T) {
	_, router, _, _ := setupTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// Unauthenticated user should be redirected to /login
	if rr.Code != http.StatusFound && rr.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect status code 302/303, got %d", rr.Code)
	}

	loc := rr.Header().Get("Location")
	if !strings.Contains(loc, "/login") {
		t.Errorf("expected redirect location /login, got %s", loc)
	}
}

func TestRoute_GuestRedirectsAuthenticatedUser(t *testing.T) {
	_, router, sm, userRepo := setupTestRouter(t)
	ctx := testhelper.SetupTestDB

	// Create user
	authService := services.NewAuthService(userRepo)
	user, _, err := authService.Register(t.Context(), "Auth User", "auth@example.com", "password123", "password123")
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}
	_ = ctx

	// Simulate session cookie
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rr := httptest.NewRecorder()

	// Put user into session
	sessionCtx, err := sm.Load(req.Context(), "")
	if err != nil {
		t.Fatalf("session load error: %v", err)
	}
	sm.Put(sessionCtx, "auth_user_id", user.ID)
	token, _, err := sm.Commit(sessionCtx)
	if err != nil {
		t.Fatalf("session commit error: %v", err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "starter_kit_session",
		Value: token,
	})

	router.ServeHTTP(rr, req)

	// Authenticated user accessing /login should be redirected to /dashboard
	if rr.Code != http.StatusFound && rr.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect status code 302/303, got %d", rr.Code)
	}

	loc := rr.Header().Get("Location")
	if !strings.Contains(loc, "/dashboard") {
		t.Errorf("expected redirect location /dashboard, got %s", loc)
	}
}

func TestRoute_LogoutDestroysSession(t *testing.T) {
	_, router, sm, userRepo := setupTestRouter(t)

	authService := services.NewAuthService(userRepo)
	user, _, err := authService.Register(t.Context(), "Logout User", "logout@example.com", "password123", "password123")
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	sessionCtx, _ := sm.Load(req.Context(), "")
	sm.Put(sessionCtx, "auth_user_id", user.ID)
	token, _, _ := sm.Commit(sessionCtx)

	req.AddCookie(&http.Cookie{
		Name:  "starter_kit_session",
		Value: token,
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusFound && rr.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect status code, got %d", rr.Code)
	}

	loc := rr.Header().Get("Location")
	if !strings.Contains(loc, "/login") {
		t.Errorf("expected redirect to /login, got %s", loc)
	}
}
