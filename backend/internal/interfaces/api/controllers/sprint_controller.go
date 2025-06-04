package controllers

import (
	"main/internal/app/queries"
	"main/internal/app/usecases"
	"main/internal/interfaces/api/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SprintController struct {
	sprintUC *usecases.SprintUseCases `di.inject:"SprintUseCases"`
}

func (c *SprintController) RegisterRoutes(router *gin.Engine) error {
	router.GET("/sprint", c.GetSprints)
	router.GET("/sprint/:id", c.GetSprint)
	router.POST("/sprint", c.CreateSprint)
	router.PUT("/sprint/:id", c.UpdateSprint)
	router.DELETE("/sprint/:id", c.DeleteSprint)
	router.POST("/sprint/:id/complete", c.CompleteSprint)
	router.POST("/sprint/:id/run", c.RunSprint)
	return nil
}

func (c *SprintController) GetSprints(ctx *gin.Context) {
	var request dto.GetSprintsRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	sprints, err := c.sprintUC.GetSprints(&queries.SprintFilter{
		TeamID:    request.Team,
		QuarterId: request.Quarter,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeSprints(*sprints))
}

func (c *SprintController) GetSprint(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	sprint, err := c.sprintUC.GetSprint(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeSprint(sprint))
}

func (c *SprintController) CreateSprint(ctx *gin.Context) {
	var request dto.CreateSprintRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if err := request.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	startDate, _ := request.GetStartDate()
	endDate, _ := request.GetEndDate()

	sprint, err := c.sprintUC.CreateSprint(&usecases.CreateSprintRequest{
		Title:        request.Title,
		StartDate:    startDate,
		EndDate:      endDate,
		TeamId:       request.TeamId,
		HoursPerUser: request.HoursPerUser,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeSprint(sprint))
}

func (c *SprintController) UpdateSprint(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	var request dto.UpdateSprintRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if err := request.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	startDate, _ := request.GetStartDate()
	endDate, _ := request.GetEndDate()

	sprint, err := c.sprintUC.UpdateSprint(&usecases.UpdateSprintRequest{
		Id:           uint(id),
		Title:        request.Title,
		StartDate:    startDate,
		EndDate:      endDate,
		TeamId:       request.TeamId,
		HoursPerUser: request.HoursPerUser,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeSprint(sprint))
}

func (c *SprintController) DeleteSprint(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	err = c.sprintUC.DeleteSprint(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (c *SprintController) CompleteSprint(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	sprint, err := c.sprintUC.CompleteSprint(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeSprint(sprint))
}

func (c *SprintController) RunSprint(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	sprint, err := c.sprintUC.RunSprint(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeSprint(sprint))
}
