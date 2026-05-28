package cmd

import (
	"main/config"
	"main/internal/bootstrap"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/goioc/di"
)

type ApiApplication struct {
	config *config.Config
	router *gin.Engine
}

func NewApiApplication() *ApiApplication {
	cfg := di.GetInstance("config").(*config.Config)
	cfg.Print()

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.AllowOrigins
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	router.Use(cors.New(corsConfig))
	if cfg.SentryDSN != "" {
		router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
	}

	app := ApiApplication{
		config: cfg,
		router: router,
	}
	app.registerDeps()
	return &app
}

func (a *ApiApplication) registerDeps() {

}

func (app *ApiApplication) Run() {
	if err := bootstrap.InitRouter(app.router); err != nil {
		panic(err)
	}

	err := app.router.Run(app.config.GetHostPort())
	if err != nil {
		panic(err)
	}
}
