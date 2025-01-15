package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port string `env:"PORT" envDefault:":8080"`
	}
	MongoDB struct {
		URI      string `env:"MONGODB_URI" envDefault:"mongodb://localhost:27017"`
		Database string `env:"MONGODB_DATABASE" envDefault:"musicdb"`
	}
	Redis struct {
		URI string `env:"REDIS_URI" envDefault:"localhost:6379"`
	}
	JWT struct {
		Secret string `env:"JWT_SECRET" envDefault:"your-secret-key"`
	}
}

func LoadConfig() (*Config, error) {
	// Cargar archivo .env
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: archivo .env no encontrado: %v", err)
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error al parsear configuración: %v", err)
	}

	// Imprimir configuración cargada (para debugging)
	log.Printf("Configuración cargada: MongoDB URI=%s, Database=%s, Redis URI=%s",
		cfg.MongoDB.URI,
		cfg.MongoDB.Database,
		cfg.Redis.URI)

	// Validar configuración mínima
	if cfg.MongoDB.URI == "" {
		return nil, fmt.Errorf("MONGODB_URI es requerido")
	}

	if cfg.MongoDB.Database == "" {
		return nil, fmt.Errorf("MONGODB_DATABASE es requerido")
	}

	return cfg, nil
}
