package handlers

import (
	"log"
	"net/http"

	gonertia "github.com/romsar/gonertia/v3"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/repositories"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/services"
)

type DashboardHandler struct {
	i              *gonertia.Inertia
	productService services.ProductService
	userRepo       repositories.UserRepository
}

func NewDashboardHandler(i *gonertia.Inertia, productService services.ProductService, userRepo repositories.UserRepository) *DashboardHandler {
	return &DashboardHandler{
		i:              i,
		productService: productService,
		userRepo:       userRepo,
	}
}

func (h *DashboardHandler) Index(w http.ResponseWriter, r *http.Request) {
	stats, err := h.productService.GetDashboardStats(r.Context(), h.userRepo)
	if err != nil {
		log.Printf("error getting dashboard stats: %v", err)
		stats = &services.DashboardStats{}
	}

	err = h.i.Render(w, r, "Dashboard", gonertia.Props{
		"stats": stats,
	})
	if err != nil {
		log.Printf("render dashboard error: %v", err)
	}
}
