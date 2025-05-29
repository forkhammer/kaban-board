package controllers

import (
	"main/internal/app/project_usecases"
	"main/internal/interfaces/api/dto"
	"main/internal/interfaces/api/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	projectUC *project_usecases.ProjectUseCases `di.inject:"ProjectUseCases"`
}

func (c *ProjectController) RegisterRoutes(router *gin.Engine) error {
	projectRoutes := router.Group("/")
	projectRoutes.Use(middleware.AuthRequiredMiddleware())
	projectRoutes.GET("/projects", c.getProjects)
	projectRoutes.POST("/projects/:id/set_team", c.setProjectTeam)
	return nil
}

func (c *ProjectController) getProjects(ctx *gin.Context) {
	projects, err := c.projectUC.GetProjects()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeProjects(projects))
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

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeProject(project))
}
