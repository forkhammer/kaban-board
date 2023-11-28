package kanban

import (
	"github.com/gin-gonic/gin"
	"main/account"
	"main/tools"
	"net/http"
	"strconv"
)

type KanbanController struct {
	userService    UserService
	columnService  ColumnService
	teamService    TeamService
	labelService   LabelService
	projectService ProjectService
}

func (c *KanbanController) RegisterRoutes(engine *gin.Engine) {
	engine.GET("/kanban-users", c.getKanbanUsers)
	engine.GET("/columns", c.getColumns)
	engine.GET("/columns/:id", c.getColumnById)
	engine.GET("/teams", c.getTeams)
	engine.GET("/teams/:id", c.getTeamById)
	engine.GET("/labels", c.getLabels)

	columnRoutes := engine.Group("/")
	columnRoutes.Use(account.AuthRequiredMiddleware())
	columnRoutes.POST("/columns", c.addColumn)
	columnRoutes.PUT("/columns/:id", c.updateColumnById)
	columnRoutes.DELETE("/columns/:id", c.deleteColumn)

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
}

func (c *KanbanController) getKanbanUsers(ctx *gin.Context) {
	k := NewKanban(tools.MemoryCacheInstance)
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
