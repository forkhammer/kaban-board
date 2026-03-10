package controllers

import (
	"main/internal/app/usecases"
	"main/internal/interfaces/api/dto"
	"main/internal/interfaces/api/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TeamController struct {
	teamUC *usecases.TeamUseCases `di.inject:"TeamUseCases"`
}

func (c *TeamController) RegisterRoutes(router gin.IRouter) error {
	router.GET("/teams", c.getTeams)
	router.GET("/teams/:id", c.getTeamById)

	protectedtRoutes := router.Group("/")
	protectedtRoutes.Use(middleware.AdminRequiredMiddleware())
	protectedtRoutes.POST("/teams", c.addTeam)
	protectedtRoutes.PUT("/teams/:id", c.updateTeam)
	protectedtRoutes.DELETE("/teams/:id", c.deleteTeam)
	return nil
}

func (c *TeamController) getTeams(ctx *gin.Context) {
	teams, err := c.teamUC.GetTeams()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeTeams(teams))
}

func (c *TeamController) getTeamById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	team, err := c.teamUC.GetTeam(uint(id))

	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeTeam(team))
}

func (c *TeamController) updateTeam(ctx *gin.Context) {
	var request dto.UpdateTeamRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	team, err := c.teamUC.Update(&usecases.UpdateTeamRequest{
		Id:     uint(id),
		Title:  request.Title,
		Groups: request.Groups,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeTeam(team))
}

func (c *TeamController) addTeam(ctx *gin.Context) {
	var request dto.CreateTeamRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	team, err := c.teamUC.Create(&usecases.CreateTeamRequest{
		Title:  request.Title,
		Groups: request.Groups,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.SerializeTeam(team))
}

func (c *TeamController) deleteTeam(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	err = c.teamUC.Delete(uint(id))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
