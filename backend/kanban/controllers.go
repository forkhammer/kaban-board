package kanban

import (
	"main/account"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goioc/di"
)

type KanbanController struct {
	userService           *UserService           `di.inject:"userService"`
	projectService        *ProjectService        `di.inject:"projectService"`
	groupService          *GroupService          `di.inject:"groupService"`
	clientSettingsService *ClientSettingsService `di.inject:"clientSettingsService"`
	kanbanSettings        *KanbanSettings        `di.inject:"kanbanSettings"`
}

func (c *KanbanController) RegisterRoutes(engine *gin.Engine) {
	engine.GET("/kanban-users", c.getKanbanUsers)
	engine.GET("/settings", c.getSettings)

	projectRoutes := engine.Group("/")
	projectRoutes.Use(account.AuthRequiredMiddleware())
	projectRoutes.GET("/projects", c.getProjects)
	projectRoutes.POST("/projects/:id/set_team", c.setProjectTeam)

	kanbanSettingsRoutes := engine.Group("/")
	kanbanSettingsRoutes.Use(account.AuthRequiredMiddleware())
	kanbanSettingsRoutes.GET("/kanban-settings", c.getKanbanSettings)
	kanbanSettingsRoutes.POST("/kanban-settings/task-type-labels", c.saveTaskTypeLabels)
}

func (c *KanbanController) getKanbanUsers(ctx *gin.Context) {
	k := di.GetInstance("kanban").(*Kanban)
	users, updateTime, err := k.GetUsers()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users, "updateTime": updateTime})
}

func (c *KanbanController) getProjects(ctx *gin.Context) {
	projects, err := c.projectService.GetProjects()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, projects)
}

func (c *KanbanController) setProjectTeam(ctx *gin.Context) {
	var request SetTeamRequest

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := c.projectService.SetTeam(uint(id), request.TeamId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, project)
}

func (c *KanbanController) getSettings(ctx *gin.Context) {
	settings := c.clientSettingsService.GetSettings()
	ctx.JSON(http.StatusOK, settings)
}

func (c *KanbanController) getKanbanSettings(ctx *gin.Context) {
	settings := c.kanbanSettings
	ctx.JSON(http.StatusOK, settings)
}

func (c *KanbanController) saveTaskTypeLabels(ctx *gin.Context) {
	var request SaveTaskTypeLabelsRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.kanbanSettings.SetTaskTypeLabels(request.Labels); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
