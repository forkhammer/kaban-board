package cmd

import (
	"context"
	"fmt"
	"main/config"
	"main/internal/app/usecases"
	"time"

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
	app.startSync(ctx)
}

func (app *WorkerApplication) startSync(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Minute * time.Duration(config.Settings.GitlabSyncPeriodMin))
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
	}
}
