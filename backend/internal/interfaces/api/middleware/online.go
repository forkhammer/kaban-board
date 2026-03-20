package middleware

import (
	"fmt"
	"main/config"
	domain "main/internal/domain/models"
	"main/internal/infra/cache"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goioc/di"
)

const onlineKeyPrefix = "online:user:"

func OnlineTrackingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		account, exists := ctx.Get("account")
		if exists && account != (*domain.Account)(nil) {
			acc := account.(*domain.Account)
			c := di.GetInstance("cache").(cache.Cache)
			ttl := time.Duration(config.Settings.OnlineUserTTLMin) * time.Minute
			key := fmt.Sprintf("%s%d", onlineKeyPrefix, acc.Id)
			c.Set(key, acc.Id, ttl)
		}
		ctx.Next()
	}
}
