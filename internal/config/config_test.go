package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Configurar variables de entorno para prueba
	os.Setenv("SERVER_PORT", ":9090")
	os.Setenv("MONGODB_URI", "mongodb://testhost:27017")
	os.Setenv("JWT_SECRET", "test-secret")

	cfg, err := LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, ":9090", cfg.Server.Port)
	assert.Equal(t, "mongodb://testhost:27017", cfg.MongoDB.URI)
	assert.Equal(t, "test-secret", cfg.JWT.Secret)

	// Limpiar variables de entorno
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("MONGODB_URI")
	os.Unsetenv("JWT_SECRET")
}

func TestLoadConfig_Defaults(t *testing.T) {
	cfg, err := LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, ":8080", cfg.Server.Port)
	assert.Equal(t, "mongodb://localhost:27017", cfg.MongoDB.URI)
}
