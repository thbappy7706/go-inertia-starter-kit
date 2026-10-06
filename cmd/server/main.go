package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	gonertia "github.com/romsar/gonertia/v3"
	"github.com/thbappy7706/go-inertia-starter-kit/database/seeders"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/config"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/database"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/handlers"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/middleware"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/repositories"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/services"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/session"
	"github.com/thbappy7706/go-inertia-starter-kit/routes"
)

func main() {
	log.Println("Starting Go + Inertia application...")

	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// 2. Connect to PostgreSQL
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// 3. Run database migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// 4. Seed initial data
	if err := seeders.Seed(db); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	// 5. Session manager and FlashProvider
	sm := session.NewManager(cfg)
	flashProvider := session.NewFlashProvider(sm.SessionManager)

	// 6. Initialize Gonertia with Vite
	i, err := gonertia.NewFromFile(
		"resources/views/app.html",
		gonertia.WithFlashProvider(flashProvider),
		gonertia.WithVersionFromFile("public/build/manifest.json"),
	)
	if err != nil {
		log.Fatalf("failed to initialize Gonertia: %v", err)
	}

	viteApp, err := gonertia.NewVite(i,
		gonertia.WithEntryPoints("frontend/main.tsx"),
		gonertia.WithBuildManifest("public/build/manifest.json"),
		gonertia.WithFallbackManifest("public/build/.vite/manifest.json"),
		gonertia.WithHotFile("public/hot"),
		gonertia.WithBuildDir("/build/"),
	)
	if err != nil {
		log.Fatalf("failed to initialize Vite instance: %v", err)
	}

	// Override viteReactRefresh with a dynamic version so the React Fast Refresh
	// preamble is always injected correctly regardless of Go vs Vite startup order.
	// gonertia captures the hot-reload state once at init; this closure re-reads the
	// hot file on every template render instead.
	if err := viteApp.Inertia.ShareTemplateFunc("viteReactRefresh", dynamicReactRefresh); err != nil {
		log.Printf("warning: could not override viteReactRefresh: %v", err)
	}

	// 7. Repositories
	userRepo := repositories.NewUserRepository(db)
	productRepo := repositories.NewProductRepository(db)

	// 8. Services
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	productService := services.NewProductService(productRepo)

	// 9. Middlewares
	authMiddleware := middleware.NewAuthMiddleware(viteApp.Inertia, sm.SessionManager, userRepo)

	// 10. Handlers
	authHandler := handlers.NewAuthHandler(viteApp.Inertia, sm.SessionManager, authService)
	dashboardHandler := handlers.NewDashboardHandler(viteApp.Inertia, productService, userRepo)
	userHandler := handlers.NewUserHandler(viteApp.Inertia, userService)
	productHandler := handlers.NewProductHandler(viteApp.Inertia, sm.SessionManager, productService)

	// 11. Routes
	router := routes.Setup(routes.RouterParams{
		Inertia:          viteApp.Inertia,
		SessionManager:   sm.SessionManager,
		AuthMiddleware:   authMiddleware,
		AuthHandler:      authHandler,
		DashboardHandler: dashboardHandler,
		UserHandler:      userHandler,
		ProductHandler:   productHandler,
	})

	// 12. HTTP Server with graceful shutdown
	serverAddr := fmt.Sprintf(":%s", cfg.AppPort)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server listening on http://localhost:%s\n", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped.")
}

// dynamicReactRefresh is a template function that re-reads public/hot on every
// render. This ensures the @vitejs/plugin-react preamble script is emitted even
// when the Go server started before `npm run dev` created the hot file.
func dynamicReactRefresh() template.HTML {
	content, err := os.ReadFile("public/hot")
	if err != nil {
		// Production mode – no hot file, no preamble needed.
		return template.HTML("")
	}

	rawURL := strings.TrimSpace(string(content))

	// Normalise to a protocol-relative URL so it works on both http and https.
	url := rawURL
	switch {
	case strings.HasPrefix(rawURL, "http://"):
		url = "//" + rawURL[7:]
	case strings.HasPrefix(rawURL, "https://"):
		url = "//" + rawURL[8:]
	}

	return template.HTML(`<script type="module">
    import RefreshRuntime from "` + url + `/@react-refresh"
    RefreshRuntime.injectIntoGlobalHook(window)
    window.$RefreshReg$ = () => {}
    window.$RefreshSig$ = () => (type) => type
    window.__vite_plugin_react_preamble_installed__ = true
</script>`)
}