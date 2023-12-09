package account

import (
	"main/tools"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/goioc/di"
)

type AccountModule struct{}

func (m *AccountModule) Init(engine *gin.Engine, connection tools.ConnectionInterface, repositoryFactory tools.RepositoryFactory) error {
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
