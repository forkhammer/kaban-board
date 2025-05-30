package middleware

import (
	"main/internal/app/usecases"
	"main/internal/interfaces/api/utils"

	"github.com/gin-gonic/gin"
	"github.com/goioc/di"
)

func JwtMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountUC := di.GetInstance("AccountUseCases").(*usecases.AccountUseCases)
		token := utils.GetTokenFromRequest(ctx)
		account, _ := accountUC.GetActiveUser(token)

		ctx.Set("account", account)

		ctx.Next()
	}
}
