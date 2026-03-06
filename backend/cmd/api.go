package cmd

import (
	"main/config"
	"main/internal/interfaces/api"
	"main/internal/interfaces/api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/goioc/di"
)

type ApiApplication struct {
	router *gin.Engine
}

func NewApiApplication() *ApiApplication {
	config.Settings.Print()

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.Settings.AllowOrigins
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	router.Use(cors.New(corsConfig))
	if config.Settings.SentryDSN != "" {
		router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
	}
	router.Use(middleware.JwtMiddleware())

	app := ApiApplication{
		router: router,
	}
	app.registerDeps()
	return &app
}

func (a *ApiApplication) registerDeps() {

}

func (app *ApiApplication) Run() {
	if err := app.initRouter(); err != nil {
		panic(err)
	}

	err := app.router.Run(config.Settings.GetHostPort())
	if err != nil {
		panic(err)
	}
}

func (app *ApiApplication) initRouter() error {
	controllers := []api.Controller{
		di.GetInstance("AccountController").(api.Controller),
		di.GetInstance("KanbanController").(api.Controller),
		di.GetInstance("ReportsController").(api.Controller),
		di.GetInstance("HealthController").(api.Controller),
		di.GetInstance("ColumnController").(api.Controller),
		di.GetInstance("TeamController").(api.Controller),
		di.GetInstance("LabelController").(api.Controller),
		di.GetInstance("UserController").(api.Controller),
		di.GetInstance("GroupController").(api.Controller),
		di.GetInstance("ProjectController").(api.Controller),
		di.GetInstance("SettingsController").(api.Controller),
		di.GetInstance("SprintController").(api.Controller),
		di.GetInstance("IssueController").(api.Controller),
		di.GetInstance("IssueBindingController").(api.Controller),
		di.GetInstance("ReleaseController").(api.Controller),
		di.GetInstance("EpicController").(api.Controller),
	}
	for _, controller := range controllers {
		err := controller.RegisterRoutes(app.router)
		if err != nil {
			return err
		}
	}
	return nil
}
