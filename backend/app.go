package main

import (
	"main/account"
	"main/config"
	"main/db"
	"main/kanban"
	"main/repository"
	"main/tools"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/goioc/di"
)

type Application struct {
	engine            *gin.Engine
	modules           []tools.AppModule
	connection        tools.ConnectionInterface
	repositoryFactory tools.RepositoryFactory
}

func NewApplication() Application {
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.Settings.AllowOrigins
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization", "ngrok-skip-browser-warning")
	router.Use(cors.New(corsConfig))
	router.Use(account.JwtMiddleware())

	modules := []tools.AppModule{
		&account.AccountModule{},
		&kanban.KanbanModule{},
	}
	connection, err := db.GetConnectionByType(config.Settings.DbType, config.Settings)
	di.RegisterBeanInstance("connection", connection)

	if err != nil {
		panic(err)
	}

	repositoryFactory, err := repository.GetRepositoryFactory(config.Settings.DbType, connection)

	if err != nil {
		panic(err)
	}

	di.RegisterBeanInstance("repositoryFactory", repositoryFactory)
	di.RegisterBeanInstance("accountRepository", repositoryFactory.GetAccountRepository())
	di.RegisterBeanInstance("columnRepository", repositoryFactory.GetColumnRepository())
	di.RegisterBeanInstance("labelRepository", repositoryFactory.GetLabelRepository())
	di.RegisterBeanInstance("projectRepository", repositoryFactory.GetProjectRepository())
	di.RegisterBeanInstance("teamRepository", repositoryFactory.GetTeamRepository())
	di.RegisterBeanInstance("userRepository", repositoryFactory.GetUserRepository())

	return Application{
		engine:            router,
		modules:           modules,
		connection:        connection,
		repositoryFactory: repositoryFactory,
	}
}

func (app *Application) Run() {
	for _, m := range app.modules {
		err := m.Init(app.engine, app.connection, app.repositoryFactory)

		if err != nil {
			panic(err)
		}
	}

	di.InitializeContainer()

	for _, m := range app.modules {
		m.RegisterRoutes(app.engine)
	}

	kb := di.GetInstance("kanban").(*kanban.Kanban)
	go kb.RunUpdater()

	app.engine.Run(config.Settings.GetHostPort())
}
