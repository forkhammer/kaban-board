package account

import (
	"main/repository"
	"main/tools"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/goioc/di"
)

type AccountModule struct{}

func (m *AccountModule) Init(engine *gin.Engine, connection repository.ConnectionInterface, repositoryFactory repository.RepositoryFactory) error {
	err := connection.Migrate(&Account{})

	if err != nil {
		return err
	}

	di.RegisterBean("accountService", reflect.TypeOf((*AccountService)(nil)))
	di.RegisterBean("accountController", reflect.TypeOf((*AccountController)(nil)))
	return nil
}

func (m *AccountModule) RegisterRoutes(engine *gin.Engine) {
	for _, c := range m.GetControllers() {
		c.RegisterRoutes(engine)
	}
}

func (m *AccountModule) GetControllers() []tools.Controller {
	return []tools.Controller{di.GetInstance("accountController").(*AccountController)}
}
