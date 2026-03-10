package utils

import (
	domain "main/internal/domain/models"

	"github.com/gin-gonic/gin"
)

func IsAdminAccount(ctx *gin.Context) bool {
	value, exists := ctx.Get("account")
	account, ok := value.(*domain.Account)
	if !exists || !ok || account == nil {
		return false
	}
	if account.Role != domain.AccountRoleAdmin {
		return false
	}
	return true
}
