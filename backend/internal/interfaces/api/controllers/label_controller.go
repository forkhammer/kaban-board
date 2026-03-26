package controllers

import (
	"main/internal/app/usecases"
	domain "main/internal/domain/models"
	"main/internal/interfaces/api/dto"
	"main/internal/interfaces/api/middleware"
	"main/internal/interfaces/api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LabelController struct {
	labelUC *usecases.LabelUseCases `di.inject:"LabelUseCases"`
}

func (c *LabelController) RegisterRoutes(router gin.IRouter) error {
	router.GET("/labels", c.getLabels)

	protectedRoutes := router.Group("/")
	protectedRoutes.Use(middleware.AdminRequiredMiddleware())
	protectedRoutes.PUT("/labels/:id", c.updateLabel)
	return nil
}

func (c *LabelController) getLabels(ctx *gin.Context) {
	labels, err := c.labelUC.GetLabels()

	if utils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeKanbanLabels(labels))
}

func (c *LabelController) updateLabel(ctx *gin.Context) {
	var request dto.UpdateLabelRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	id := ctx.Param("id")

	err := c.labelUC.Update(&usecases.UpdateLabelRequest{
		Title:         string(id),
		AltName:       request.AltName,
		BindingStatus: (*domain.IssueBindingStatus)(request.BindingStatus),
		Priority:      (*domain.IssueBindingPriority)(request.Priority),
	})

	if utils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
