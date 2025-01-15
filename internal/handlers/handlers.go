package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/JonayMedina/go-search-api/internal/services"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	musicService services.MusicServiceInterface
	userService  *services.UserService
}

func NewHandlers(musicService services.MusicServiceInterface, userService *services.UserService) *Handlers {
	return &Handlers{
		musicService: musicService,
		userService:  userService,
	}
}

func (handlers *Handlers) Search(ctx *gin.Context) {
	query := strings.TrimSpace(ctx.Query("q"))
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "El parámetro de búsqueda no puede estar vacío",
		})
		return
	}

	if len(query) < 3 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "La búsqueda debe tener al menos 3 caracteres",
		})
		return
	}

	songs, err := handlers.musicService.SearchMusic(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, songs)
}

func (handlers *Handlers) SearchHistory(ctx *gin.Context) {
	// Aquí deberías implementar la lógica para obtener el historial de búsquedas
	// Por ahora, devolvemos un mensaje temporal
	ctx.JSON(http.StatusOK, gin.H{
		"message": "search history - to be implemented",
	})
}

func (handlers *Handlers) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Unix(),
	})
}

func (handlers *Handlers) GetUsers(ctx *gin.Context) {
	log.Println("Recibida solicitud para obtener usuarios")
	users, err := handlers.userService.GetUsers(ctx.Request.Context())
	if err != nil {
		log.Printf("Error obteniendo usuarios: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios"})
		return
	}
	log.Printf("Enviando respuesta con %d usuarios", len(users))
	ctx.JSON(http.StatusOK, users)
}

func (handlers *Handlers) GetUserByUsername(ctx *gin.Context) {
	log.Println("Recibida solicitud para obtener usuarios")
	username := ctx.Param("username")
	user, err := handlers.userService.GetUserByUsername(ctx.Request.Context(), username)
	if err != nil {
		log.Printf("Error obteniendo usuario: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuario"})
		return
	}
	log.Printf("Enviando respuesta con usuario: %+v", user)
	ctx.JSON(http.StatusOK, user)

}

func validateRequest(ctx *gin.Context) error {
	if ctx.ContentType() != "application/json" {
		return fmt.Errorf("content-type debe ser application/json")
	}

	if ctx.Request.Body == nil {
		return fmt.Errorf("el cuerpo de la petición no puede estar vacío")
	}

	return nil
}
