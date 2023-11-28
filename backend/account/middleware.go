package account

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func JwtMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		account, _ := GetCurrentAccountFromContext(ctx)

		ctx.Set("account", account)

		ctx.Next()
	}
}

func AuthRequiredMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		account, exists := ctx.Get("account")
		if account == nil || !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentification required"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
