package services

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/thbappy7706/go-inertia-starter-kit/internal/models"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/repositories"
	"gorm.io/gorm"
)

var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type CreateProductDTO struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	Price       float64 `json:"price"`
	Status      bool    `json:"status"`
}

type UpdateProductDTO struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	Price       float64 `json:"price"`
	Status      bool    `json:"status"`
}

type DashboardStats struct {
	TotalUsers     int64                    `json:"total_users"`
	TotalProducts  int64                    `json:"total_products"`
	ActiveProducts int64                    `json:"active_products"`
	RecentProducts []models.ProductResponse `json:"recent_products"`
}

type ProductService interface {
	List(ctx context.Context, search string, status *bool, page, perPage int) (models.PaginatedResponse[models.ProductResponse], error)
	FindByID(ctx context.Context, id uint) (*models.ProductResponse, error)
	Create(ctx context.Context, dto CreateProductDTO) (*models.ProductResponse, map[string]string, error)
	Update(ctx context.Context, id uint, dto UpdateProductDTO) (*models.ProductResponse, map[string]string, error)
	Delete(ctx context.Context, id uint) error
	GetDashboardStats(ctx context.Context, userRepo repositories.UserRepository) (*DashboardStats, error)
}

type productService struct {
	productRepo repositories.ProductRepository
}

func NewProductService(productRepo repositories.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) List(ctx context.Context, search string, status *bool, page, perPage int) (models.PaginatedResponse[models.ProductResponse], error) {
	products, total, err := s.productRepo.List(ctx, search, status, page, perPage)
	if err != nil {
		return models.PaginatedResponse[models.ProductResponse]{}, err
	}

	responses := make([]models.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = p.ToResponse()
	}

	return models.NewPaginatedResponse(responses, total, page, perPage), nil
}

func (s *productService) FindByID(ctx context.Context, id uint) (*models.ProductResponse, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := product.ToResponse()
	return &res, nil
}

func (s *productService) Create(ctx context.Context, dto CreateProductDTO) (*models.ProductResponse, map[string]string, error) {
	valErrors := s.validate(ctx, 0, dto.Name, dto.Slug, dto.Price)
	if len(valErrors) > 0 {
		return nil, valErrors, nil
	}

	var desc *string
	if dto.Description != nil && strings.TrimSpace(*dto.Description) != "" {
		trimmed := strings.TrimSpace(*dto.Description)
		desc = &trimmed
	}

	product := &models.Product{
		Name:        strings.TrimSpace(dto.Name),
		Slug:        strings.ToLower(strings.TrimSpace(dto.Slug)),
		Description: desc,
		Price:       dto.Price,
		Status:      dto.Status,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, nil, err
	}

	res := product.ToResponse()
	return &res, nil, nil
}

func (s *productService) Update(ctx context.Context, id uint, dto UpdateProductDTO) (*models.ProductResponse, map[string]string, error) {
	existing, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	valErrors := s.validate(ctx, id, dto.Name, dto.Slug, dto.Price)
	if len(valErrors) > 0 {
		return nil, valErrors, nil
	}

	var desc *string
	if dto.Description != nil && strings.TrimSpace(*dto.Description) != "" {
		trimmed := strings.TrimSpace(*dto.Description)
		desc = &trimmed
	}

	existing.Name = strings.TrimSpace(dto.Name)
	existing.Slug = strings.ToLower(strings.TrimSpace(dto.Slug))
	existing.Description = desc
	existing.Price = dto.Price
	existing.Status = dto.Status

	if err := s.productRepo.Update(ctx, existing); err != nil {
		return nil, nil, err
	}

	res := existing.ToResponse()
	return &res, nil, nil
}

func (s *productService) Delete(ctx context.Context, id uint) error {
	return s.productRepo.Delete(ctx, id)
}

func (s *productService) validate(ctx context.Context, currentID uint, name, slug string, price float64) map[string]string {
	valErrors := make(map[string]string)

	name = strings.TrimSpace(name)
	slug = strings.ToLower(strings.TrimSpace(slug))

	if len(name) < 2 {
		valErrors["name"] = "Name must be at least 2 characters."
	} else if len(name) > 255 {
		valErrors["name"] = "Name cannot exceed 255 characters."
	}

	if slug == "" {
		valErrors["slug"] = "Slug is required."
	} else if !slugRegex.MatchString(slug) {
		valErrors["slug"] = "Slug must contain only lowercase letters, numbers, and hyphens."
	} else {
		existing, err := s.productRepo.FindBySlug(ctx, slug)
		if err == nil && existing != nil && existing.ID != currentID {
			valErrors["slug"] = "This slug is already taken."
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			valErrors["slug"] = "Error verifying slug uniqueness."
		}
	}

	if price <= 0 {
		valErrors["price"] = "Price must be a positive number."
	}

	return valErrors
}

func (s *productService) GetDashboardStats(ctx context.Context, userRepo repositories.UserRepository) (*DashboardStats, error) {
	totalUsers, err := userRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	totalProducts, err := s.productRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	activeProducts, err := s.productRepo.CountActive(ctx)
	if err != nil {
		return nil, err
	}

	recent, err := s.productRepo.Recent(ctx, 5)
	if err != nil {
		return nil, err
	}

	recentResponses := make([]models.ProductResponse, len(recent))
	for i, p := range recent {
		recentResponses[i] = p.ToResponse()
	}

	return &DashboardStats{
		TotalUsers:     totalUsers,
		TotalProducts:  totalProducts,
		ActiveProducts: activeProducts,
		RecentProducts: recentResponses,
	}, nil
}
