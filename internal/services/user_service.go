package services

import (
	"context"

	"github.com/thbappy7706/go-inertia-starter-kit/internal/models"
	"github.com/thbappy7706/go-inertia-starter-kit/internal/repositories"
)

type UserService interface {
	List(ctx context.Context, search string, page, perPage int) (models.PaginatedResponse[models.UserResponse], error)
	FindByID(ctx context.Context, id uint) (*models.UserResponse, error)
	Count(ctx context.Context) (int64, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) List(ctx context.Context, search string, page, perPage int) (models.PaginatedResponse[models.UserResponse], error) {
	users, total, err := s.userRepo.List(ctx, search, page, perPage)
	if err != nil {
		return models.PaginatedResponse[models.UserResponse]{}, err
	}

	responses := make([]models.UserResponse, len(users))
	for i, u := range users {
		responses[i] = u.ToResponse()
	}

	return models.NewPaginatedResponse(responses, total, page, perPage), nil
}

func (s *userService) FindByID(ctx context.Context, id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := user.ToResponse()
	return &res, nil
}

func (s *userService) Count(ctx context.Context) (int64, error) {
	return s.userRepo.Count(ctx)
}
