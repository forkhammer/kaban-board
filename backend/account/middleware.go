package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goioc/di"
)

func JwtMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountService := di.GetInstance("accountService").(*AccountService)

		account, _ := accountService.GetCurrentAccountFromContext(ctx)

		ctx.Set("account", account)

		ctx.Next()
	}
}

func AuthRequiredMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, exists := ctx.Get("account")
		account := data.(*Account)
		if account == nil || !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentification required"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
