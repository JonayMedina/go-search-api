package bootstrap

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/JonayMedina/go-search-api/internal/config"
	"github.com/JonayMedina/go-search-api/internal/db/mongodb"
	"github.com/JonayMedina/go-search-api/internal/services"
)

type Services struct {
	MusicService *services.MusicService
	UserService  *services.UserService
}

func InitServices(cfg *config.Config) (*Services, error) {
	mongoClient, err := initMongoDB(cfg)
	if err != nil {
		return nil, err
	}

	redisClient := initRedis(cfg)

	musicService := services.NewMusicService(mongoClient, redisClient)
	userRepo := mongodb.NewUserRepository(mongoClient, cfg.MongoDB.Database)
	userService := services.NewUserService(userRepo)

	return &Services{
		MusicService: musicService,
		UserService:  userService,
	}, nil
}

func initMongoDB(cfg *config.Config) (*mongo.Client, error) {
	ctx := context.Background()
	log.Printf("Conectando a MongoDB en: %s", cfg.MongoDB.URI)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		return nil, fmt.Errorf("error conectando a MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("error haciendo ping a MongoDB: %v", err)
	}

	log.Println("Conexión exitosa a MongoDB")
	return client, nil
}

func initRedis(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.Redis.URI,
	})
}
