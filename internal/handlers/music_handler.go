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

func (musicHandler *MusicHandler) Search(ctx *gin.Context) {
	query := ctx.Query("q")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	songs, err := musicHandler.service.SearchMusic(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, songs)
}
