package kanban

import (
	"main/cache"
	"main/repository"
	"main/repository/models"
	"main/tools"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/goioc/di"
)

type KanbanModule struct{}

func (m *KanbanModule) Init(engine *gin.Engine, connection repository.ConnectionInterface, repositoryFactory repository.RepositoryFactory) error {
	err := connection.Migrate(&models.Project{}, &models.User{}, &models.Group{}, &models.Team{}, &models.Label{}, &models.Column{}, &models.KVElement{})

	if err != nil {
		return err
	}

	di.RegisterBean("userService", reflect.TypeOf((*UserService)(nil)))
	di.RegisterBean("labelService", reflect.TypeOf((*LabelService)(nil)))
	di.RegisterBean("projectService", reflect.TypeOf((*ProjectService)(nil)))
	di.RegisterBean("teamService", reflect.TypeOf((*TeamService)(nil)))
	di.RegisterBean("columnService", reflect.TypeOf((*ColumnService)(nil)))
	di.RegisterBean("groupService", reflect.TypeOf((*GroupService)(nil)))
	di.RegisterBeanInstance("kanban", NewKanban(cache.MemoryCacheInstance))
	di.RegisterBean("kanbanController", reflect.TypeOf((*KanbanController)(nil)))
	di.RegisterBean("clientSettingsService", reflect.TypeOf((*ClientSettingsService)(nil)))
	di.RegisterBean("kvStore", reflect.TypeOf((*KVStore)(nil)))
	di.RegisterBean("kanbanSettings", reflect.TypeOf((*KanbanSettings)(nil)))
	return nil
}

func (m *KanbanModule) RegisterRoutes(engine *gin.Engine) {
	for _, c := range m.GetControllers() {
		c.RegisterRoutes(engine)
	}
}

func (m *KanbanModule) GetControllers() []tools.Controller {
	return []tools.Controller{di.GetInstance("kanbanController").(*KanbanController)}
}
