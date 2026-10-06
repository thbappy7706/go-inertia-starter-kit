package handlers

import (
	"log"
	"net/http"

	gonertia "github.com/romsar/gonertia/v3"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/services"
)

type UserHandler struct {
	i           *gonertia.Inertia
	userService services.UserService
}

func NewUserHandler(i *gonertia.Inertia, userService services.UserService) *UserHandler {
	return &UserHandler{
		i:           i,
		userService: userService,
	}
}

func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	page := getIntQuery(r, "page", 1)
	perPage := 10

	users, err := h.userService.List(r.Context(), search, page, perPage)
	if err != nil {
		log.Printf("list users error: %v", err)
	}

	err = h.i.Render(w, r, "Users/Index", gonertia.Props{
		"users": users,
		"filters": map[string]any{
			"search": search,
		},
	})
	if err != nil {
		log.Printf("render users error: %v", err)
	}
}
