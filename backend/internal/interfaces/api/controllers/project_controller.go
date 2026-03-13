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

type ProjectController struct {
	projectUC *usecases.ProjectUseCases `di.inject:"ProjectUseCases"`
}

func (c *ProjectController) RegisterRoutes(router gin.IRouter) error {
	router.GET("/projects", c.getProjects)
	router.GET("/projects/:id", c.getProject)

	projectRoutes := router.Group("/")
	projectRoutes.Use(middleware.AdminRequiredMiddleware())
	projectRoutes.POST("/projects/:id/set_team", c.setProjectTeam)
	return nil
}

func (c *ProjectController) getProjects(ctx *gin.Context) {
	var request dto.GetProjectsRequest

	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	projects, err := c.projectUC.GetProjects(&queries.ProjectFilter{
		TeamID: request.TeamId,
		Search: request.Search,
	})

	if utils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeProjects(projects))
}

func (c *ProjectController) getProject(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	project, err := c.projectUC.GetProject(uint(id))

	if utils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeProject(project))
}

func (c *ProjectController) setProjectTeam(ctx *gin.Context) {
	var request dto.SetTeamRequest

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if err = ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	project, err := c.projectUC.SetTeam(uint(id), request.TeamId)

	if utils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeProject(project))
}
