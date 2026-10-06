package repositories

import (
	"context"
	"strings"

	"github.com/thbappy7706/go-inertia-starter-kit/internal/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindByID(ctx context.Context, id uint) (*models.Product, error)
	FindBySlug(ctx context.Context, slug string) (*models.Product, error)
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, search string, status *bool, page, perPage int) ([]models.Product, int64, error)
	Count(ctx context.Context) (int64, error)
	CountActive(ctx context.Context) (int64, error)
	Recent(ctx context.Context, limit int) ([]models.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindByID(ctx context.Context, id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindBySlug(ctx context.Context, slug string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Where("slug = ?", strings.TrimSpace(slug)).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

func (r *productRepository) List(ctx context.Context, search string, status *bool, page, perPage int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Product{})

	if search = strings.TrimSpace(search); search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name ILIKE ? OR slug ILIKE ? OR description ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := query.Order("id DESC").Offset(offset).Limit(perPage).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Product{}).Count(&count).Error
	return count, err
}

func (r *productRepository) CountActive(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Product{}).Where("status = ?", true).Count(&count).Error
	return count, err
}

func (r *productRepository) Recent(ctx context.Context, limit int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&products).Error
	return products, err
}
