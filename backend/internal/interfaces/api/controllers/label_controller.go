package controllers

import (
	"main/internal/app/label_usecases"
	"main/internal/interfaces/api/dto"
	"main/internal/interfaces/api/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LabelController struct {
	labelUC *label_usecases.LabelUseCases `di.inject:"LabelUseCases"`
}

func (c *LabelController) RegisterRoutes(router *gin.Engine) error {
	router.GET("/labels", c.getLabels)

	protectedRoutes := router.Group("/")
	protectedRoutes.Use(middleware.AuthRequiredMiddleware())
	protectedRoutes.PUT("/labels/:id", c.updateLabel)
	return nil
}

func (c *LabelController) getLabels(ctx *gin.Context) {
	labels, err := c.labelUC.GetLabels()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeLabels(labels))
}

func (c *LabelController) updateLabel(ctx *gin.Context) {
	var request dto.UpdateLabelRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	id := ctx.Param("id")

	err := c.labelUC.Update(&label_usecases.UpdateLabelRequest{
		Title:   string(id),
		AltName: request.AltName,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
