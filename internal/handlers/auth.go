package handlers

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/JonayMedina/go-search-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (handlers *Handlers) Login(ctx *gin.Context) {
	if err := validateRequest(ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Usuario y contraseña son requeridos",
		})
		return
	}

	// Aquí deberías validar las credenciales contra tu base de datos
	log.Printf("Intentando autenticar usuario: %s", req.Username)

	user, err := handlers.userService.GetUserByUsername(ctx.Request.Context(), req.Username)
	if err != nil {

		log.Printf("Error al obtener usuario: %v", err)
		if err.Error() == "mongo: no documents in result" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuario"})
		return
	}

	log.Printf("Usuario obtenido: %+v", user)

	valid := user.ComparePassword(req.Password)
	if !valid {
		log.Println("Credenciales inválidas")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	// Por ahora, simulamos una validación básica
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	ctx.JSON(http.StatusOK, models.LoginResponse{
		Token: tokenString,
		User:  *user,
	})
}

func (handlers *Handlers) Register(ctx *gin.Context) {
	if err := validateRequest(ctx); err != nil {
		log.Printf("Error en validación de request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req models.RegisterRequest
	body, _ := io.ReadAll(ctx.Request.Body)
	log.Printf("Request body: %s", string(body))

	// Importante: restaurar el body para que pueda ser leído nuevamente
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error en binding JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, models.RegisterResponse{
			Message: "Datos inválidos",
			Details: map[string]string{
				"username": "Debe tener entre 3 y 30 caracteres alfanuméricos",
				"password": "Debe tener al menos 6 caracteres",
				"email":    "Debe ser un email válido",
			},
		})
		return
	}

	log.Printf("Request procesado: %+v", req)

	if strings.TrimSpace(req.Username) == "" ||
		strings.TrimSpace(req.Password) == "" ||
		strings.TrimSpace(req.Email) == "" {
		ctx.JSON(http.StatusBadRequest, models.RegisterResponse{
			Message: "Los campos email, username y password son requeridos",
		})
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := handlers.userService.Register(ctx.Request.Context(), user); err != nil {
		if err.Error() == "username or email already exists" {
			ctx.JSON(http.StatusConflict, models.RegisterResponse{
				Message: "El usuario o email ya existe",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.RegisterResponse{
			Message: "Error al registrar usuario",
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.RegisterResponse{
		Message: "Usuario registrado exitosamente",
		User:    *user,
	})
}
