package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/JonayMedina/go-search-api/internal/bootstrap"
	"github.com/JonayMedina/go-search-api/internal/config"
	"github.com/JonayMedina/go-search-api/internal/handlers"
	"github.com/JonayMedina/go-search-api/internal/router"
)

func main() {
	setupLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("No se pudo cargar la configuración:", err)
	}

	services, err := bootstrap.InitServices(cfg)
	if err != nil {
		log.Fatal("Error inicializando servicios:", err)
	}

	handlers := handlers.NewHandlers(services.MusicService, services.UserService)
	r := router.SetupRouter(handlers)

	log.Printf("Servidor iniciado en el puerto %s", cfg.Server.Port)
	if err := r.Run(cfg.Server.Port); err != nil {
		log.Fatal("Error iniciando el servidor:", err)
	}
}

func setupLogger() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	gin.DefaultWriter = os.Stdout
}
