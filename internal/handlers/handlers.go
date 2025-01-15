package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/JonayMedina/go-search-api/internal/models"
	"github.com/JonayMedina/go-search-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func (h *Handlers) Search(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "El parámetro de búsqueda no puede estar vacío",
		})
		return
	}

	if len(query) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "La búsqueda debe tener al menos 3 caracteres",
		})
		return
	}

	songs, err := h.musicService.SearchMusic(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, songs)
}

func (h *Handlers) Login(c *gin.Context) {
	if err := validateRequest(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos inválidos",
			"details": map[string]string{
				"username": "Campo requerido",
				"password": "Campo requerido",
			},
			"message": err.Error(),
		})
		return
	}

	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Usuario y contraseña son requeridos",
		})
		return
	}

	// Aquí deberías validar las credenciales contra tu base de datos
	// Por ahora, simulamos una validación básica
	if req.Username == "admin" && req.Password == "password" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": req.Username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString([]byte("your-secret-key"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
}

func (h *Handlers) Register(c *gin.Context) {
	if err := validateRequest(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("Register", c.Request.Body)
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos inválidos",
			"details": map[string]string{
				"username": "Debe tener entre 3 y 30 caracteres alfanuméricos",
				"password": "Debe tener al menos 6 caracteres",
				"email":    "Debe ser un email válido",
			},
			"message": err.Error(),
		})
		return
	}

	if strings.TrimSpace(req.Username) == "" ||
		strings.TrimSpace(req.Password) == "" ||
		strings.TrimSpace(req.Email) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Los campos no pueden estar vacíos",
		})
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.userService.Register(c.Request.Context(), user); err != nil {
		if err.Error() == "username or email already exists" {
			c.JSON(http.StatusConflict, gin.H{
				"error": "El usuario o email ya existe",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al registrar usuario",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuario registrado exitosamente",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func (h *Handlers) SearchHistory(c *gin.Context) {
	// Aquí deberías implementar la lógica para obtener el historial de búsquedas
	// Por ahora, devolvemos un mensaje temporal
	c.JSON(http.StatusOK, gin.H{
		"message": "search history - to be implemented",
	})
}

func (h *Handlers) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		// Remover el prefijo "Bearer " si existe
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Unix(),
	})
}

func validateRequest(c *gin.Context) error {
	if c.ContentType() != "application/json" {
		return fmt.Errorf("content-type debe ser application/json")
	}

	if c.Request.Body == nil {
		return fmt.Errorf("el cuerpo de la petición no puede estar vacío")
	}

	return nil
}
