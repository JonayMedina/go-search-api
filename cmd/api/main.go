package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/JonayMedina/go-search-api/internal/handlers"
	"github.com/JonayMedina/go-search-api/internal/services"
)

func main() {
	// Inicializar MongoDB
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
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
	musicService := services.NewMusicService(mongoClient, redisClient)
	musicHandler := handlers.NewMusicHandler(musicService)

	// Configurar router
	r := gin.Default()

	// Rutas públicas
	r.POST("/auth/login", handlers.Login)

	// Rutas protegidas
	protected := r.Group("/api")
	protected.Use(handlers.AuthMiddleware())
	{
		protected.GET("/search", musicHandler.Search)
	}

	r.Run(":8080")
}
