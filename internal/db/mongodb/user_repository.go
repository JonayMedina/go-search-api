package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/JonayMedina/go-search-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	client *mongo.Client
	db     string
}

func NewUserRepository(client *mongo.Client, db string) *UserRepository {
	return &UserRepository{
		client: client,
		db:     db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	collection := r.client.Database(r.db).Collection("users")

	// Verificar si el usuario ya existe
	exists, err := r.exists(ctx, user.Username, user.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("username or email already exists")
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := user.HashPassword(); err != nil {
		return err
	}

	_, err = collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) exists(ctx context.Context, username, email string) (bool, error) {
	collection := r.client.Database(r.db).Collection("users")

	filter := bson.M{
		"$or": []bson.M{
			{"username": username},
			{"email": email},
		},
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
