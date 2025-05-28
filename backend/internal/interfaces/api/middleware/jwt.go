package middleware

import (
	"main/internal/app/account_usecases"
	"main/internal/interfaces/api/utils"

	"github.com/gin-gonic/gin"
	"github.com/goioc/di"
)

func JwtMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		activeUser := di.GetInstance("ActiveUserUseCase").(*account_usecases.ActiveUserUseCase)
		token := utils.GetTokenFromRequest(ctx)
		account, _ := activeUser.Execute(token)

		ctx.Set("account", account)

		ctx.Next()
	}
}
