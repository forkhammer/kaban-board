package controllers

import (
	"main/internal/app/usecases"
	"main/internal/interfaces/api/dto"
	"main/internal/interfaces/api/middleware"
	"main/internal/interfaces/api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SettingsController struct {
	settingsUC *usecases.SettingsUseCases `di.inject:"SettingsUseCases"`
}

func (c *SettingsController) RegisterRoutes(router gin.IRouter) error {
	router.GET("/settings", c.getSettings)

	protectedRoutes := router.Group("/")
	protectedRoutes.Use(middleware.AdminRequiredMiddleware())
	protectedRoutes.GET("/kanban-settings", c.getKanbanSettings)
	protectedRoutes.POST("/kanban-settings/task-type-labels", c.saveTaskTypeLabels)
	return nil
}

func (c *SettingsController) getSettings(ctx *gin.Context) {
	settings := c.settingsUC.GetClientSettings()
	ctx.JSON(http.StatusOK, settings)
}

func (c *SettingsController) getKanbanSettings(ctx *gin.Context) {
	settings, err := c.settingsUC.GetKanbantSettings()
	if utils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, settings)
}

func (c *SettingsController) saveTaskTypeLabels(ctx *gin.Context) {
	var request dto.SaveTaskTypeLabelsRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if _, err := c.settingsUC.SetTaskTypeLabels(request.Labels); utils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
