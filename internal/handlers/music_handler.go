package handlers

import (
	"net/http"

	"github.com/JonayMedina/go-search-api/internal/services"
	"github.com/gin-gonic/gin"
)

type MusicHandler struct {
	service *services.MusicService
}

func NewMusicHandler(service *services.MusicService) *MusicHandler {
	return &MusicHandler{service: service}
}

func (h *MusicHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	songs, err := h.service.SearchMusic(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, songs)
}
