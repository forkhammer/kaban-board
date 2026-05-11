package main

import (
	"context"
	"fmt"
	"main/cmd"
	"main/config"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"
	"time"

	sentry "github.com/getsentry/sentry-go"
	"github.com/goioc/di"
)

type Application struct {
	api    *cmd.ApiApplication
	worker *cmd.WorkerApplication
}

func NewApplication() *Application {
	cfg := config.GetConfig()
	app := &Application{}
	app.Init(cfg)
	return app
}

func initSentry(cfg *config.Config) {
	dsn := cfg.SentryDSN
	if dsn == "" {
		return
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      cfg.SentryEnvironment,
		TracesSampleRate: cfg.SentrySampleRate,
	})
	if err != nil {
		fmt.Printf("Sentry init failed: %v\n", err)
		return
	}
	fmt.Println("Sentry initialized")
}

func (app *Application) Init(cfg *config.Config) {
	initSentry(cfg)
	initDI()
}

func (app *Application) Run() {
	defer sentry.Flush(2 * time.Second)

	err := app.migrate()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	app.api = cmd.NewApiApplication()
	app.worker = cmd.NewWorkerApplication()

	app.worker.Run(ctx)
	app.api.Run()
}

func (app *Application) migrate() error {
	connection := di.GetInstance("connection").(interfaces.ConnectionInterface)
	return connection.Migrate(models.ALL_MODELS...)
}

func main() {
	app := NewApplication()
	app.Run()
}
