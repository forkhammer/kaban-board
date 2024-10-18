package tools

import (
	"main/repository"

	"github.com/gin-gonic/gin"
)

type AppModule interface {
	Init(engine *gin.Engine, connection repository.ConnectionInterface, repositoryFactory repository.RepositoryFactory) error
	RegisterRoutes(engine *gin.Engine)
	GetControllers() []Controller
}
