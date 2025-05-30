package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	domain "main/internal/domain/models"
	"main/internal/interfaces/api/dto"
)

func AuthRequiredMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		account, exists := ctx.Get("account")
		if account == (*domain.Account)(nil) || !exists {
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Authentification required"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
