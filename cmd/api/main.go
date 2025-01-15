package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/JonayMedina/go-search-api/internal/config"
	"github.com/JonayMedina/go-search-api/internal/cors"
	"github.com/JonayMedina/go-search-api/internal/db/mongodb"
	"github.com/JonayMedina/go-search-api/internal/handlers"
	"github.com/JonayMedina/go-search-api/internal/services"
)

func main() {
	// Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("No se pudo cargar la configuración:", err)
	}

	// Inicializar MongoDB
	mongoClient, err := initMongoDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Inicializar Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URI"),
	})
	defer redisClient.Close()

	// Inicializar servicios
	musicService := initMusicService(mongoClient, redisClient)
	userRepo := mongodb.NewUserRepository(mongoClient, cfg.MongoDB.Database)
	userService := services.NewUserService(userRepo)
	handlers := handlers.NewHandlers(musicService, userService)

	// Configurar router
	router := setupRouter(handlers)

	log.Printf("Servidor iniciado en el puerto %s", cfg.Server.Port)
	router.Run(cfg.Server.Port)
}

func setupRouter(handlers *handlers.Handlers) *gin.Engine {
	router := gin.Default()

	// Middleware global
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(cors.Default())

	// Rutas de salud
	router.GET("/health", handlers.HealthCheck)

	// Rutas de autenticación
	auth := router.Group("/auth")
	{
		auth.POST("/login", handlers.Login)
		auth.POST("/register", handlers.Register)
	}

	// Rutas protegidas
	api := router.Group("/api")
	api.Use(handlers.AuthMiddleware())
	{
		api.GET("/search", handlers.Search)
		api.GET("/history", handlers.SearchHistory)
	}

	return router
}

func initMongoDB(cfg *config.Config) (*mongo.Client, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		return nil, fmt.Errorf("error conectando a MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("error haciendo ping a MongoDB: %v", err)
	}

	return client, nil
}

func initMusicService(mongoClient *mongo.Client, redisClient *redis.Client) *services.MusicService {
	return services.NewMusicService(mongoClient, redisClient)
}
