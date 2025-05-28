package controllers

import (
	"main/internal/app/column_usecases"
	"main/internal/interfaces/api/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ColumnController struct {
	listUC *column_usecases.ListColumnsUseCase `di.inject:"ListColumnsUseCase"`
}

func (c *ColumnController) RegisterRoutes(router *gin.Engine) error {
	router.GET("/columns", c.getColumns)
	// router.GET("/columns/:id", c.getColumnById)

	// columnRoutes := router.Group("/")
	// columnRoutes.Use(middleware.AuthRequiredMiddleware())
	// columnRoutes.POST("/columns", c.addColumn)
	// columnRoutes.PUT("/columns/:id", c.updateColumnById)
	// columnRoutes.DELETE("/columns/:id", c.deleteColumn)
	// columnRoutes.POST("/columns/save_ordering", c.saveColumnOrdering)
	return nil
}

func (c *ColumnController) getColumns(ctx *gin.Context) {
	columns, err := c.listUC.Execute()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeColumns(columns))
}
