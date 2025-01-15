package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JonayMedina/go-search-api/internal/models"
	"github.com/JonayMedina/go-search-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlers_HealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h := &Handlers{}
	h.HealthCheck(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
	assert.NotNil(t, response["time"])
}

func TestHandlers_Search(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockMusicService := &services.MockMusicService{}
	mockUserService := &services.UserService{}
	mockMusicService.On("SearchMusic", mock.Anything, "test").Return([]models.Song{}, nil)

	h := NewHandlers(mockMusicService, mockUserService)

	// Configurar la solicitud de prueba
	c.Request, _ = http.NewRequest(http.MethodGet, "/api/search?q=test", nil)
	c.Query("q") // Simular parámetro de consulta

	h.Search(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
