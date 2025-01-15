package services

import (
	"context"

	"github.com/JonayMedina/go-search-api/internal/db/mongodb"
	"github.com/JonayMedina/go-search-api/internal/models"
)

type UserService struct {
	userRepo *mongodb.UserRepository
}

func NewUserService(userRepo *mongodb.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (service *UserService) Register(ctx context.Context, user *models.User) error {
	return service.userRepo.Create(ctx, user)
}
