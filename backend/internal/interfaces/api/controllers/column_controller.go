package controllers

import (
	"main/internal/app/column_usecases"
	"main/internal/interfaces/api/dto"
	"main/internal/interfaces/api/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ColumnController struct {
	listUC     *column_usecases.ListColumnsUseCase    `di.inject:"ListColumnsUseCase"`
	retrieveUC *column_usecases.RetrieveColumnUseCase `di.inject:"RetrieveColumnUseCase"`
	updateUC   *column_usecases.UpdateColumnUseCase   `di.inject:"UpdateColumnUseCase"`
	createUC   *column_usecases.CreateColumnUseCase   `di.inject:"CreateColumnUseCase"`
}

func (c *ColumnController) RegisterRoutes(router *gin.Engine) error {
	router.GET("/columns", c.getColumns)
	router.GET("/columns/:id", c.getColumnById)

	protectedRoutes := router.Group("/")
	protectedRoutes.Use(middleware.AuthRequiredMiddleware())
	protectedRoutes.POST("/columns", c.addColumn)
	protectedRoutes.PUT("/columns/:id", c.updateColumn)
	// protectedRoutes.DELETE("/columns/:id", c.deleteColumn)
	// protectedRoutes.POST("/columns/save_ordering", c.saveColumnOrdering)
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

func (c *ColumnController) getColumnById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	column, err := c.retrieveUC.Execute(uint(id))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeColumn(column))
}

func (c *ColumnController) addColumn(ctx *gin.Context) {
	var request dto.CreateColumnRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	column, err := c.createUC.Execute(&column_usecases.CreateColumnRequest{
		Name:   request.Name,
		Labels: request.Labels,
		TeamId: request.TeamId,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.SerializeColumn(column))
}

func (c *ColumnController) updateColumn(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	var request dto.UpdateColumnRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	column, err := c.updateUC.Execute(&column_usecases.UpdateColumnRequest{
		Id:     uint(id),
		Name:   request.Name,
		Labels: request.Labels,
		TeamId: request.TeamId,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.SerializeColumn(column))
}
