package web

import (
	"main/config"
	"main/internal/interfaces/api"
	"main/internal/interfaces/api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/goioc/di"
)

type Application struct {
	router *gin.Engine
}

func NewApplication() *Application {
	config.Settings.Print()

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.Settings.AllowOrigins
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	router.Use(cors.New(corsConfig))
	router.Use(middleware.JwtMiddleware())

	app := Application{
		router: router,
	}
	app.registerDeps()
	return &app
}

func (a *Application) registerDeps() {

}

func (app *Application) Run() {
	if err := app.initRouter(); err != nil {
		panic(err)
	}

	app.router.Run(config.Settings.GetHostPort())
}

func (app *Application) initRouter() error {
	controllers := []api.Controller{
		di.GetInstance("AccountController").(api.Controller),
		di.GetInstance("BoardController").(api.Controller),
		di.GetInstance("ReportsController").(api.Controller),
		di.GetInstance("HealthController").(api.Controller),
		di.GetInstance("ColumnController").(api.Controller),
		di.GetInstance("TeamController").(api.Controller),
		di.GetInstance("LabelController").(api.Controller),
		di.GetInstance("UserController").(api.Controller),
		di.GetInstance("GroupController").(api.Controller),
		di.GetInstance("ProjectController").(api.Controller),
		di.GetInstance("SettingsController").(api.Controller),
	}
	for _, controller := range controllers {
		err := controller.RegisterRoutes(app.router)
		if err != nil {
			return err
		}
	}
	return nil
}
