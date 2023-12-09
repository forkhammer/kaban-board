package tools

import (
	"github.com/gin-gonic/gin"
)

type AppModule interface {
	Init(engine *gin.Engine, connection ConnectionInterface, repositoryFactory RepositoryFactory) error
	RegisterRoutes(engine *gin.Engine)
	GetControllers() []Controller
}
