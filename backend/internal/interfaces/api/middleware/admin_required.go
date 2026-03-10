package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	domain "main/internal/domain/models"
	"main/internal/interfaces/api/dto"
)

func AdminRequiredMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value, exists := ctx.Get("account")
		account, ok := value.(*domain.Account)
		if !exists || !ok || account == nil {
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Authentification required"})
			ctx.Abort()
			return
		}
		if account.Role != domain.AccountRoleAdmin {
			ctx.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "Admin role required"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
