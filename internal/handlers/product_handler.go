package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	gonertia "github.com/romsar/gonertia/v3"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/services"
)

type ProductHandler struct {
	i              *gonertia.Inertia
	sm             *scs.SessionManager
	productService services.ProductService
}

func NewProductHandler(i *gonertia.Inertia, sm *scs.SessionManager, productService services.ProductService) *ProductHandler {
	return &ProductHandler{
		i:              i,
		sm:             sm,
		productService: productService,
	}
}

func (h *ProductHandler) Index(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	statusParam := r.URL.Query().Get("status")
	page := getIntQuery(r, "page", 1)
	perPage := 10

	var status *bool
	if statusParam == "active" {
		b := true
		status = &b
	} else if statusParam == "inactive" {
		b := false
		status = &b
	}

	products, err := h.productService.List(r.Context(), search, status, page, perPage)
	if err != nil {
		log.Printf("list products error: %v", err)
	}

	err = h.i.Render(w, r, "Products/Index", gonertia.Props{
		"products": products,
		"filters": map[string]any{
			"search": search,
			"status": statusParam,
		},
	})
	if err != nil {
		log.Printf("render products error: %v", err)
	}
}

func (h *ProductHandler) ShowCreate(w http.ResponseWriter, r *http.Request) {
	if err := h.i.Render(w, r, "Products/Create"); err != nil {
		log.Printf("render create product error: %v", err)
	}
}

func (h *ProductHandler) Store(w http.ResponseWriter, r *http.Request) {
	var dto services.CreateProductDTO
	if err := parseJSONOrForm(r, &dto); err != nil {
		h.sm.Put(r.Context(), "flash_error", "Invalid request body.")
		h.i.Back(w, r)
		return
	}

	_, valErrors, err := h.productService.Create(r.Context(), dto)
	if err != nil {
		log.Printf("create product error: %v", err)
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

	h.sm.Put(r.Context(), "flash_success", "Product created successfully.")
	h.i.Redirect(w, r, "/products")
}

func (h *ProductHandler) ShowEdit(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.sm.Put(r.Context(), "flash_error", "Invalid product ID.")
		h.i.Redirect(w, r, "/products")
		return
	}

	product, err := h.productService.FindByID(r.Context(), uint(id))
	if err != nil {
		h.sm.Put(r.Context(), "flash_error", "Product not found.")
		h.i.Redirect(w, r, "/products")
		return
	}

	if err := h.i.Render(w, r, "Products/Edit", gonertia.Props{
		"product": product,
	}); err != nil {
		log.Printf("render edit product error: %v", err)
	}
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.sm.Put(r.Context(), "flash_error", "Invalid product ID.")
		h.i.Back(w, r)
		return
	}

	var dto services.UpdateProductDTO
	if err := parseJSONOrForm(r, &dto); err != nil {
		h.sm.Put(r.Context(), "flash_error", "Invalid request body.")
		h.i.Back(w, r)
		return
	}

	_, valErrors, err := h.productService.Update(r.Context(), uint(id), dto)
	if err != nil {
		log.Printf("update product error: %v", err)
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

	h.sm.Put(r.Context(), "flash_success", "Product updated successfully.")
	h.i.Redirect(w, r, "/products")
}

func (h *ProductHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.sm.Put(r.Context(), "flash_error", "Invalid product ID.")
		h.i.Redirect(w, r, "/products")
		return
	}

	if err := h.productService.Delete(r.Context(), uint(id)); err != nil {
		log.Printf("delete product error: %v", err)
		h.sm.Put(r.Context(), "flash_error", "Failed to delete product.")
		h.i.Redirect(w, r, "/products")
		return
	}

	h.sm.Put(r.Context(), "flash_success", "Product deleted successfully.")
	h.i.Redirect(w, r, "/products")
}