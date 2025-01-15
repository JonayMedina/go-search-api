package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Server struct {
		Port string `env:"SERVER_PORT" envDefault:":8080"`
	}
	MongoDB struct {
		URI      string `env:"MONGODB_URI" envDefault:"mongodb://localhost:27017"`
		Database string `env:"MONGODB_DB" envDefault:"musicdb"`
	}
	Redis struct {
		URI      string `env:"REDIS_URI" envDefault:"localhost:6379"`
		Password string `env:"REDIS_PASSWORD" envDefault:""`
	}
	JWT struct {
		Secret string `env:"JWT_SECRET" envDefault:"your-secret-key"`
	}
}

func loadEnv(cfg *Config) error {
	return env.Parse(cfg)
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := loadEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
