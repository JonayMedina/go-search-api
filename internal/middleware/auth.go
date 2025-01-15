package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Printf("Verificando autenticación para ruta: %s", ctx.Request.URL.Path)
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			log.Println("Token no proporcionado")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required!!"})
			ctx.Abort()
			return
		}

		// Remover el prefijo "Bearer " si existe
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		log.Println("Token validado exitosamente")
		ctx.Next()
	}
}
