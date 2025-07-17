package controllers

import (
	"main/internal/app/queries"
	"main/internal/app/usecases"
	"main/internal/interfaces/api/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EpicController struct {
	epicUC *usecases.EpicUseCases `di.inject:"EpicUseCases"`
}

func (c *EpicController) RegisterRoutes(router *gin.Engine) error {
	routes := router.Group("/")
	routes.GET("/epic", c.getEpics)
	routes.GET("/epic/:id", c.getEpic)
	return nil
}

func (c *EpicController) getEpics(ctx *gin.Context) {
	var request dto.GetEpicsRequest

	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	epics, err := c.epicUC.GetEpics(&queries.EpicFilter{
		Search:    request.Search,
		ProjectId: request.Project,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeEpics(epics))
}

func (c *EpicController) getEpic(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	epic, err := c.epicUC.GetEpic(uint(id))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeEpic(epic))
}
