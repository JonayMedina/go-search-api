package services

import (
	"context"
	"log"

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
	log.Println("Register from user service", user)
	return service.userRepo.Create(ctx, user)
}

func (service *UserService) GetUsers(ctx context.Context) ([]models.User, error) {
	log.Println("Obteniendo lista de usuarios...")
	users, err := service.userRepo.GetUsers(ctx)
	if err != nil {
		log.Printf("Error obteniendo usuarios: %v", err)
		return nil, err
	}
	log.Printf("Se encontraron %d usuarios", len(users))
	return users, nil
}

func (service *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	log.Printf("Obteniendo usuario por username: %s", username)
	user, err := service.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		log.Printf("Error obteniendo usuario: %v", err)
		return nil, err
	}
	return user, nil

}
