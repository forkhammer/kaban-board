package cmd

import (
	"context"
	"fmt"
	"main/config"
	"main/internal/app/usecases"
	"time"

	sentry "github.com/getsentry/sentry-go"
	"github.com/goioc/di"
)

type WorkerApplication struct {
	config *config.Config
	syncUc *usecases.SyncUseCases
}

func NewWorkerApplication() *WorkerApplication {
	return &WorkerApplication{
		config: di.GetInstance("config").(*config.Config),
		syncUc: di.GetInstance("SyncUseCases").(*usecases.SyncUseCases),
	}
}

func (app *WorkerApplication) Run(ctx context.Context) {
	if app.config.GitlabSyncEnabled {
		app.startSync(ctx)
	}
}

func (app *WorkerApplication) startSync(ctx context.Context) {
	go func() {
		defer sentry.Recover()
		ticker := time.NewTicker(time.Minute * time.Duration(app.config.GitlabSyncPeriodMin))
		defer ticker.Stop()

		app.syncIteration()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				app.syncIteration()
			}
		}
	}()
}

func (app *WorkerApplication) syncIteration() {
	fmt.Println("Sync tick")
	err := app.syncUc.Sync()

	if err != nil {
		fmt.Println(err.Error())
		sentry.CaptureException(err)
	}
}
