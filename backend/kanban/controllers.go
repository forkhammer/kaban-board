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
	columnService         *ColumnService         `di.inject:"columnService"`
	teamService           *TeamService           `di.inject:"teamService"`
	labelService          *LabelService          `di.inject:"labelService"`
	projectService        *ProjectService        `di.inject:"projectService"`
	groupService          *GroupService          `di.inject:"groupService"`
	clientSettingsService *ClientSettingsService `di.inject:"clientSettingsService"`
	kanbanSettings        *KanbanSettings        `di.inject:"kanbanSettings"`
}

func (c *KanbanController) RegisterRoutes(engine *gin.Engine) {
	engine.GET("/kanban-users", c.getKanbanUsers)
	engine.GET("/columns", c.getColumns)
	engine.GET("/columns/:id", c.getColumnById)
	engine.GET("/teams", c.getTeams)
	engine.GET("/teams/:id", c.getTeamById)
	engine.GET("/labels", c.getLabels)
	engine.GET("/settings", c.getSettings)
	engine.GET("/groups", c.getGroups)
	engine.GET("/groups/:id", c.getGroupById)

	columnRoutes := engine.Group("/")
	columnRoutes.Use(account.AuthRequiredMiddleware())
	columnRoutes.POST("/columns", c.addColumn)
	columnRoutes.PUT("/columns/:id", c.updateColumnById)
	columnRoutes.DELETE("/columns/:id", c.deleteColumn)
	columnRoutes.POST("/columns/save_ordering", c.saveColumnOrdering)

	teamRoutes := engine.Group("/")
	teamRoutes.Use(account.AuthRequiredMiddleware())
	teamRoutes.POST("/teams", c.addTeam)
	teamRoutes.PUT("/teams/:id", c.updateTeamById)
	teamRoutes.DELETE("/teams/:id", c.deleteTeam)

	userRoutes := engine.Group("/")
	userRoutes.Use(account.AuthRequiredMiddleware())
	userRoutes.GET("/users", c.getUsers)
	userRoutes.POST("/users/:id/visibility", c.setUserVisibility)

	projectRoutes := engine.Group("/")
	projectRoutes.Use(account.AuthRequiredMiddleware())
	projectRoutes.GET("/projects", c.getProjects)
	projectRoutes.POST("/projects/:id/set_team", c.setProjectTeam)

	kanbanSettingsRoutes := engine.Group("/")
	kanbanSettingsRoutes.Use(account.AuthRequiredMiddleware())
	kanbanSettingsRoutes.GET("/kanban-settings", c.getKanbanSettings)
	kanbanSettingsRoutes.POST("/kanban-settings/task-type-labels", c.saveTaskTypeLabels)

	labelRoutes := engine.Group("/")
	labelRoutes.Use(account.AuthRequiredMiddleware())
	labelRoutes.PUT("/labels/:id", c.updateLabelById)
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

func (c *KanbanController) getUsers(ctx *gin.Context) {
	users, err := c.userService.GetUsers()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (c *KanbanController) setUserVisibility(ctx *gin.Context) {
	var request SetUserVisibilityRequest

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.SetUserVisibility(int(id), request.Visible)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *KanbanController) getColumns(ctx *gin.Context) {
	columns, err := c.columnService.GetAllColumns()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, columns)
}

func (c *KanbanController) getColumnById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	column, err := c.columnService.GetColumnById(int(id))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, column)
}

func (c *KanbanController) updateColumnById(ctx *gin.Context) {
	var request UpdateColumnRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	column, err := c.columnService.UpdateColumn(int(id), &request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, column)
}

func (c *KanbanController) addColumn(ctx *gin.Context) {
	var request CreateColumnRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	column, err := c.columnService.CreateColumn(&request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusCreated, column)
}

func (c *KanbanController) deleteColumn(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.columnService.DeleteColumnById(int(id))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (c *KanbanController) saveColumnOrdering(ctx *gin.Context) {
	request := make([]SetColumnOrderRequest, 0)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	columns, err := c.columnService.SaveOrdering(request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, columns)
}

func (c *KanbanController) getTeams(ctx *gin.Context) {
	teams, err := c.teamService.GetAllTeams()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, teams)
}

func (c *KanbanController) getTeamById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := c.teamService.GetTeamById(int(id))

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, team)
}

func (c *KanbanController) updateTeamById(ctx *gin.Context) {
	var request UpdateTeamRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := c.teamService.UpdateTeam(int(id), &request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, team)
}

func (c *KanbanController) addTeam(ctx *gin.Context) {
	var request CreateTeamRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := c.teamService.CreateTeam(&request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusCreated, team)
}

func (c *KanbanController) deleteTeam(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.teamService.DeleteTeamById(int(id))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (c *KanbanController) getLabels(ctx *gin.Context) {
	labels, err := c.labelService.GetAllKanbanLabels()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, labels)
}

func (c *KanbanController) updateLabelById(ctx *gin.Context) {
	var request UpdateKanbanLabelRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")

	if err := c.labelService.UpdateKanbanLabel(id, &request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
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

func (c *KanbanController) getGroups(ctx *gin.Context) {
	groups, err := c.groupService.GetGroups()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, groups)
}

func (c *KanbanController) getGroupById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := c.groupService.GetGroupById(int(id))

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, group)
}
