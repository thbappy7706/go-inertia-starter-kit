package routes

import (
	"net/http"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	gonertia "github.com/romsar/gonertia/v3"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/handlers"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/middleware"
)

type RouterParams struct {
	Inertia          *gonertia.Inertia
	SessionManager   *scs.SessionManager
	AuthMiddleware   *middleware.AuthMiddleware
	AuthHandler      *handlers.AuthHandler
	DashboardHandler *handlers.DashboardHandler
	UserHandler      *handlers.UserHandler
	ProductHandler   *handlers.ProductHandler
}

func Setup(p RouterParams) *chi.Mux {
	r := chi.NewRouter()

	// Global standard middlewares
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// Session middleware
	r.Use(p.SessionManager.LoadAndSave)

	// Gonertia & Auth sharing middlewares
	r.Use(p.Inertia.Middleware)
	r.Use(p.AuthMiddleware.ShareAuth)

	// Static asset file server for /build/ and public files
	filesDir := http.Dir("./public")
	fileServer(r, "/build", filesDir)
	fileServer(r, "/public", filesDir)

	// Guest routes
	r.Group(func(guest chi.Router) {
		guest.Use(p.AuthMiddleware.RequireGuest)

		guest.Get("/login", p.AuthHandler.ShowLogin)
		guest.Post("/login", p.AuthHandler.Login)
		guest.Get("/register", p.AuthHandler.ShowRegister)
		guest.Post("/register", p.AuthHandler.Register)
	})

	// Protected routes
	r.Group(func(auth chi.Router) {
		auth.Use(p.AuthMiddleware.RequireAuth)

		auth.Post("/logout", p.AuthHandler.Logout)

		auth.Get("/", func(w http.ResponseWriter, r *http.Request) {
			p.Inertia.Redirect(w, r, "/dashboard")
		})
		auth.Get("/dashboard", p.DashboardHandler.Index)

		auth.Get("/users", p.UserHandler.Index)

		auth.Get("/products", p.ProductHandler.Index)
		auth.Get("/products/create", p.ProductHandler.ShowCreate)
		auth.Post("/products", p.ProductHandler.Store)
		auth.Get("/products/{id}/edit", p.ProductHandler.ShowEdit)
		auth.Put("/products/{id}", p.ProductHandler.Update)
		auth.Delete("/products/{id}", p.ProductHandler.Destroy)
	})

	return r
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
