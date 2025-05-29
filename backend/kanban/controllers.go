package kanban

import (
	"net/http"

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
