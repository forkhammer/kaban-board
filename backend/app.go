package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"main/account"
	"main/config"
	"main/kanban"
	"main/tools"
)

type Application struct {
	engine      *gin.Engine
	controllers []Controller
}

func NewApplication() Application {
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:4200", "https://0286-94-41-238-118.ngrok-free.app"}
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization", "ngrok-skip-browser-warning")
	router.Use(cors.New(corsConfig))
	router.Use(account.JwtMiddleware())

	account.MigrateModels()
	kanban.MigrateModels()

	return Application{
		engine: router,
	}
}

func (app *Application) Run() {
	app.controllers = []Controller{
		&account.AccountController{},
		&kanban.KanbanController{},
	}

	for _, c := range app.controllers {
		c.RegisterRoutes(app.engine)
	}

	kb := kanban.NewKanban(tools.MemoryCacheInstance)
	go kb.RunUpdater()

	app.engine.Run(config.Settings.GetHostPort())
}
