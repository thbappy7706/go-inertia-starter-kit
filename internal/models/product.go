package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Slug        string         `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	Description *string        `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"type:numeric(12,2);not null;default:0" json:"price"`
	Status      bool           `gorm:"default:true;not null" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProductResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description"`
	Price       float64   `json:"price"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Product) ToResponse() ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Price:       p.Price,
		Status:      p.Status,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
