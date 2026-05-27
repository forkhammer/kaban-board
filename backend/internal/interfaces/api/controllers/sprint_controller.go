package controllers

import (
	"main/internal/app/queries"
	"main/internal/app/usecases"
	"main/internal/interfaces/api/dto"
	"main/internal/interfaces/api/middleware"
	"main/internal/interfaces/api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SprintController struct {
	sprintUC              *usecases.SprintUseCases              `di.inject:"SprintUseCases"`
	sprintUserSettingsUC  *usecases.SprintUserSettingsUseCases   `di.inject:"SprintUserSettingsUseCases"`
}

func (c *SprintController) RegisterRoutes(router gin.IRouter) error {
	router.GET("/sprint", c.GetSprints)
	router.GET("/sprint/quarters", c.GetQuarters)
	router.GET("/sprint/:id", c.GetSprint)

	privateRoutes := router.Group("/")
	privateRoutes.Use(middleware.AdminRequiredMiddleware())
	privateRoutes.POST("/sprint", c.CreateSprint)
	privateRoutes.PUT("/sprint/:id", c.UpdateSprint)
	privateRoutes.DELETE("/sprint/:id", c.DeleteSprint)
	privateRoutes.POST("/sprint/:id/complete", c.CompleteSprint)
	privateRoutes.POST("/sprint/:id/run", c.RunSprint)
	privateRoutes.GET("/sprint/:id/user-settings", c.GetSprintUserSettings)
	privateRoutes.POST("/sprint/:id/user-settings", c.CreateOrUpdateSprintUserSettings)
	privateRoutes.DELETE("/sprint/:id/user-settings/:user_id", c.DeleteSprintUserSettings)
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
	if utils.HandleException(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeSprints(sprints))
}

func (c *SprintController) GetSprint(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	sprint, err := c.sprintUC.GetSprint(uint(id))
	if utils.HandleException(ctx, err) {
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
	if utils.HandleException(ctx, err) {
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
	if utils.HandleException(ctx, err) {
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
	if utils.HandleException(ctx, err) {
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
	if utils.HandleException(ctx, err) {
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
	if utils.HandleException(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeSprint(sprint))
}

func (c *SprintController) GetQuarters(ctx *gin.Context) {
	quarters, err := c.sprintUC.GetQuarters()
	if utils.HandleException(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeQuarters(quarters))
}

func (c *SprintController) GetSprintUserSettings(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	sprintId := uint(id)
	settings, err := c.sprintUserSettingsUC.GetSprintUserSettings(&queries.SprintUserSettingsFilter{
		SprintId: &sprintId,
	})
	if utils.HandleException(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeSprintUserSettingsList(settings))
}

func (c *SprintController) CreateOrUpdateSprintUserSettings(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	var request dto.CreateSprintUserSettingsRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	settings, err := c.sprintUserSettingsUC.CreateOrUpdateSprintUserSettings(uint(id), &usecases.CreateSprintUserSettingsRequest{
		UserId:       request.UserId,
		HoursPerUser: request.HoursPerUser,
	})
	if utils.HandleException(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, dto.SerializeSprintUserSettings(settings))
}

func (c *SprintController) DeleteSprintUserSettings(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	userId, err := strconv.ParseUint(ctx.Param("user_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	err = c.sprintUserSettingsUC.DeleteSprintUserSettings(uint(id), uint(userId))
	if utils.HandleException(ctx, err) {
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}
