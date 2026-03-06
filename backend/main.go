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
	app := &Application{}
	app.Init()
	return app
}

func initSentry() {
	dsn := config.Settings.SentryDSN
	if dsn == "" {
		return
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      config.Settings.SentryEnvironment,
		TracesSampleRate: config.Settings.SentrySampleRate,
	})
	if err != nil {
		fmt.Printf("Sentry init failed: %v\n", err)
		return
	}
	fmt.Println("Sentry initialized")
}

func (app *Application) Init() {
	initSentry()
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
