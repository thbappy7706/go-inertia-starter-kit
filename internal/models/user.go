package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"size:255;not null" json:"name"`
	Email           string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password        string         `gorm:"size:255;not null" json:"-"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserResponse struct {
	ID              uint       `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:              u.ID,
		Name:            u.Name,
		Email:           u.Email,
		EmailVerifiedAt: u.EmailVerifiedAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}
