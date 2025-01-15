package router

import (
	"github.com/JonayMedina/go-search-api/internal/cors"
	"github.com/JonayMedina/go-search-api/internal/handlers"
	"github.com/JonayMedina/go-search-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(handlers *handlers.Handlers) *gin.Engine {
	router := gin.Default()

	// Middleware global
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(cors.Default())

	setupHealthRoutes(router, handlers)
	setupAuthRoutes(router, handlers)
	setupAPIRoutes(router, handlers)

	return router
}

func setupHealthRoutes(router *gin.Engine, handlers *handlers.Handlers) {
	router.GET("/health", handlers.HealthCheck)
}

func setupAuthRoutes(router *gin.Engine, handlers *handlers.Handlers) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", handlers.Login)
		auth.POST("/register", handlers.Register)
	}
}

func setupAPIRoutes(router *gin.Engine, handlers *handlers.Handlers) {
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{

		api.GET("/get-users", handlers.GetUsers)
		api.GET("/search", handlers.Search)
	}
}
