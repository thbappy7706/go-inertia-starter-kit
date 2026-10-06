package middleware

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
	gonertia "github.com/romsar/gonertia/v3"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/models"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/repositories"
)

type contextKey string

const (
	UserContextKey contextKey = "current_user"
)

type AuthMiddleware struct {
	i        *gonertia.Inertia
	sm       *scs.SessionManager
	userRepo repositories.UserRepository
}

func NewAuthMiddleware(i *gonertia.Inertia, sm *scs.SessionManager, userRepo repositories.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		i:        i,
		sm:       sm,
		userRepo: userRepo,
	}
}

func (m *AuthMiddleware) ShareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var safeUser *models.UserResponse
		if userID := m.sm.Get(ctx, "auth_user_id"); userID != nil {
			if id, ok := userID.(uint); ok {
				user, err := m.userRepo.FindByID(ctx, id)
				if err == nil && user != nil {
					res := user.ToResponse()
					safeUser = &res
					ctx = context.WithValue(ctx, UserContextKey, user)
				}
			}
		}

		ctx = gonertia.SetProp(ctx, "auth", map[string]any{
			"user": safeUser,
		})

		flashSuccess := m.sm.PopString(ctx, "flash_success")
		flashError := m.sm.PopString(ctx, "flash_error")
		flash := map[string]any{}
		if flashSuccess != "" {
			flash["success"] = flashSuccess
		}
		if flashError != "" {
			flash["error"] = flashError
		}
		ctx = gonertia.SetProp(ctx, "flash", flash)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserContextKey)
		if user == nil {
			m.i.Redirect(w, r, "/login")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) RequireGuest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserContextKey)
		if user != nil {
			m.i.Redirect(w, r, "/dashboard")
			return
		}
		next.ServeHTTP(w, r)
	})
}
