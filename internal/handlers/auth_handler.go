package handlers

import (
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	gonertia "github.com/romsar/gonertia/v3"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/services"
)

type AuthHandler struct {
	i           *gonertia.Inertia
	sm          *scs.SessionManager
	authService services.AuthService
}

func NewAuthHandler(i *gonertia.Inertia, sm *scs.SessionManager, authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		i:           i,
		sm:          sm,
		authService: authService,
	}
}

func (h *AuthHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	if err := h.i.Render(w, r, "Auth/Login"); err != nil {
		log.Printf("render error: %v", err)
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := parseJSONOrForm(r, &req); err != nil {
		h.sm.Put(r.Context(), "flash_error", "Invalid request body.")
		h.i.Back(w, r)
		return
	}

	user, valErrors, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf("login error: %v", err)
		h.sm.Put(r.Context(), "flash_error", "An internal error occurred.")
		h.i.Back(w, r)
		return
	}

	if len(valErrors) > 0 {
		errs := make(gonertia.ValidationErrors)
		for k, v := range valErrors {
			errs[k] = v
		}
		setValidationErrors(r, errs)
		h.i.Back(w, r)
		return
	}

	if err := h.sm.RenewToken(r.Context()); err != nil {
		log.Printf("renew token error: %v", err)
	}

	h.sm.Put(r.Context(), "auth_user_id", user.ID)
	h.sm.Put(r.Context(), "flash_success", "Login successful.")
	h.i.Redirect(w, r, "/dashboard")
}

func (h *AuthHandler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	if err := h.i.Render(w, r, "Auth/Register"); err != nil {
		log.Printf("render error: %v", err)
	}
}

type RegisterRequest struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := parseJSONOrForm(r, &req); err != nil {
		h.sm.Put(r.Context(), "flash_error", "Invalid request body.")
		h.i.Back(w, r)
		return
	}

	user, valErrors, err := h.authService.Register(r.Context(), req.Name, req.Email, req.Password, req.PasswordConfirmation)
	if err != nil {
		log.Printf("registration error: %v", err)
		h.sm.Put(r.Context(), "flash_error", "An internal error occurred.")
		h.i.Back(w, r)
		return
	}

	if len(valErrors) > 0 {
		errs := make(gonertia.ValidationErrors)
		for k, v := range valErrors {
			errs[k] = v
		}
		setValidationErrors(r, errs)
		h.i.Back(w, r)
		return
	}

	if err := h.sm.RenewToken(r.Context()); err != nil {
		log.Printf("renew token error: %v", err)
	}

	h.sm.Put(r.Context(), "auth_user_id", user.ID)
	h.sm.Put(r.Context(), "flash_success", "Account created successfully.")
	h.i.Redirect(w, r, "/dashboard")
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.sm.Destroy(r.Context()); err != nil {
		log.Printf("destroy session error: %v", err)
	}

	h.sm.Put(r.Context(), "flash_success", "You have been logged out.")
	h.i.Redirect(w, r, "/login")
}