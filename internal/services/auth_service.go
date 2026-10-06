package services

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"github.com/thbappy7706/go-inertia-starter-kit/internal/models"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("Invalid email or password.")
)

type AuthService interface {
	Register(ctx context.Context, name, email, password, passwordConfirm string) (*models.User, map[string]string, error)
	Login(ctx context.Context, email, password string) (*models.User, map[string]string, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(ctx context.Context, name, email, password, passwordConfirm string) (*models.User, map[string]string, error) {
	valErrors := make(map[string]string)

	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if len(name) < 2 {
		valErrors["name"] = "Name must be at least 2 characters."
	} else if len(name) > 255 {
		valErrors["name"] = "Name cannot exceed 255 characters."
	}

	if email == "" {
		valErrors["email"] = "Email is required."
	} else if _, err := mail.ParseAddress(email); err != nil {
		valErrors["email"] = "Please enter a valid email address."
	} else {
		existing, err := s.userRepo.FindByEmail(ctx, email)
		if err == nil && existing != nil {
			valErrors["email"] = "This email is already registered."
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, err
		}
	}

	if len(password) < 8 {
		valErrors["password"] = "Password must be at least 8 characters."
	} else if password != passwordConfirm {
		valErrors["password_confirmation"] = "Passwords do not match."
	}

	if len(valErrors) > 0 {
		return nil, valErrors, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	user := &models.User{
		Name:     name,
		Email:    strings.ToLower(email),
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, nil, err
	}

	return user, nil, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*models.User, map[string]string, error) {
	valErrors := make(map[string]string)

	email = strings.TrimSpace(email)
	if email == "" {
		valErrors["email"] = "Email is required."
	}
	if password == "" {
		valErrors["password"] = "Password is required."
	}

	if len(valErrors) > 0 {
		return nil, valErrors, nil
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			valErrors["email"] = "Invalid email or password."
			return nil, valErrors, nil
		}
		return nil, nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		valErrors["email"] = "Invalid email or password."
		return nil, valErrors, nil
	}

	return user, nil, nil
}
