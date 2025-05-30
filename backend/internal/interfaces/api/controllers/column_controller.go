package controllers

import (
	"main/internal/app/usecases"
	domain "main/internal/domain/models"
	"main/internal/interfaces/api/dto"
	"main/internal/interfaces/api/middleware"
	"main/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ColumnController struct {
	columnUC *usecases.ColumnUseCases `di.inject:"ColumnUseCases"`
}

func (c *ColumnController) RegisterRoutes(router *gin.Engine) error {
	router.GET("/columns", c.getColumns)
	router.GET("/columns/:id", c.getColumnById)

	protectedRoutes := router.Group("/")
	protectedRoutes.Use(middleware.AuthRequiredMiddleware())
	protectedRoutes.POST("/columns", c.addColumn)
	protectedRoutes.PUT("/columns/:id", c.updateColumn)
	protectedRoutes.DELETE("/columns/:id", c.deleteColumn)
	protectedRoutes.POST("/columns/save_ordering", c.saveColumnOrdering)
	return nil
}

func (c *ColumnController) getColumns(ctx *gin.Context) {
	columns, err := c.columnUC.List()

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

	column, err := c.columnUC.Retrieve(uint(id))

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

	column, err := c.columnUC.Create(&usecases.CreateColumnRequest{
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

	column, err := c.columnUC.Update(&usecases.UpdateColumnRequest{
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

func (c *ColumnController) deleteColumn(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	err = c.columnUC.Delete(uint(id))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (c *ColumnController) saveColumnOrdering(ctx *gin.Context) {
	request := make(dto.SetColumnOrderRequest, 0)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ordering := utils.Map(request, func(o dto.SetColumnOrder) usecases.ColumnOrdering {
		return usecases.ColumnOrdering{
			Id:    o.Id,
			Order: o.Order,
		}
	})
	columns, err := c.columnUC.Ordering(ordering)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, utils.Map(*columns, func(column domain.Column) dto.ColumnDto {
		return *dto.SerializeColumn(&column)
	}))
}
