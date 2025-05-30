package controllers

import "github.com/gin-gonic/gin"

type HealthController struct {
}

func (c *HealthController) RegisterRoutes(router *gin.Engine) error {
	return nil
}
