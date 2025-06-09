package main

import (
	"context"
	"main/cmd"
	"main/config"
	app_services "main/internal/app/services"
	"main/internal/app/usecases"
	"main/internal/infra/db/implementation"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/gitlab"
	"main/internal/infra/persistance/models"
	"main/internal/infra/persistance/repo"
	column_spec "main/internal/infra/persistance/spec/column"
	group_spec "main/internal/infra/persistance/spec/group"
	issue_spec "main/internal/infra/persistance/spec/issue"
	label_spec "main/internal/infra/persistance/spec/label"
	sprint_spec "main/internal/infra/persistance/spec/sprint"
	user_spec "main/internal/infra/persistance/spec/user"
	"main/internal/infra/services"
	"main/internal/interfaces/api/controllers"
	"reflect"

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

func (app *Application) Init() {
	connection, err := implementation.GetConnectionByType(interfaces.DbType(config.Settings.DbType), config.Settings)
	di.RegisterBeanInstance("connection", connection)

	if err != nil {
		panic(err)
	}

	di.RegisterBeanInstance("config", config.Settings)
	di.RegisterBeanInstance("db", connection)
	di.RegisterBeanFactory("gitlab", di.Singleton, func(ctx context.Context) (interface{}, error) {
		return gitlab.NewGitlabClient(config.Settings.GitlabUrl, config.Settings.GitlabToken), nil
	})

	di.RegisterBean("AccountRepository", reflect.TypeOf((*repo.AccountRepository)(nil)))
	di.RegisterBean("ColumnRepository", reflect.TypeOf((*repo.ColumnRepository)(nil)))
	di.RegisterBean("GroupRepository", reflect.TypeOf((*repo.GroupRepository)(nil)))
	di.RegisterBean("IssueRepository", reflect.TypeOf((*repo.IssueRepository)(nil)))
	di.RegisterBean("KeyValueRepository", reflect.TypeOf((*repo.KeyValueRepository)(nil)))
	di.RegisterBean("LabelRepository", reflect.TypeOf((*repo.LabelRepository)(nil)))
	di.RegisterBean("ProjectRepository", reflect.TypeOf((*repo.ProjectRepository)(nil)))
	di.RegisterBean("ReleaseRepository", reflect.TypeOf((*repo.ReleaseRepository)(nil)))
	di.RegisterBean("TeamRepository", reflect.TypeOf((*repo.TeamRepository)(nil)))
	di.RegisterBean("UserRepository", reflect.TypeOf((*repo.UserRepository)(nil)))
	di.RegisterBean("SettingsRepository", reflect.TypeOf((*repo.SettingsRepository)(nil)))
	di.RegisterBean("SprintRepository", reflect.TypeOf((*repo.SprintRepository)(nil)))
	di.RegisterBean("IssueBindingRepository", reflect.TypeOf((*repo.IssueBindingRepository)(nil)))

	di.RegisterBean("LabelQuery", reflect.TypeOf((*label_spec.LabelQueryImpl)(nil)))
	di.RegisterBean("ColumnQuery", reflect.TypeOf((*column_spec.ColumnQueryImpl)(nil)))
	di.RegisterBean("GroupQuery", reflect.TypeOf((*group_spec.GroupQueryImpl)(nil)))
	di.RegisterBean("UserQuery", reflect.TypeOf((*user_spec.UserQueryImpl)(nil)))
	di.RegisterBean("SprintQuery", reflect.TypeOf((*sprint_spec.SprintQueryImpl)(nil)))
	di.RegisterBean("IssueQuery", reflect.TypeOf((*issue_spec.IssueQueryImpl)(nil)))

	di.RegisterBean("JWTService", reflect.TypeOf((*services.JWTService)(nil)))
	di.RegisterBean("PasswordService", reflect.TypeOf((*services.PasswordService)(nil)))
	di.RegisterBean("LabelService", reflect.TypeOf((*app_services.LabelService)(nil)))

	di.RegisterBean("AccountUseCases", reflect.TypeOf((*usecases.AccountUseCases)(nil)))
	di.RegisterBean("ColumnUseCases", reflect.TypeOf((*usecases.ColumnUseCases)(nil)))
	di.RegisterBean("TeamUseCases", reflect.TypeOf((*usecases.TeamUseCases)(nil)))
	di.RegisterBean("LabelUseCases", reflect.TypeOf((*usecases.LabelUseCases)(nil)))
	di.RegisterBean("UserUseCases", reflect.TypeOf((*usecases.UserUseCases)(nil)))
	di.RegisterBean("GroupUseCases", reflect.TypeOf((*usecases.GroupUseCases)(nil)))
	di.RegisterBean("ProjectUseCases", reflect.TypeOf((*usecases.ProjectUseCases)(nil)))
	di.RegisterBean("SettingsUseCases", reflect.TypeOf((*usecases.SettingsUseCases)(nil)))
	di.RegisterBean("SyncUseCases", reflect.TypeOf((*usecases.SyncUseCases)(nil)))
	di.RegisterBean("KanbanUseCases", reflect.TypeOf((*usecases.KanbanUseCases)(nil)))
	di.RegisterBean("SprintUseCases", reflect.TypeOf((*usecases.SprintUseCases)(nil)))
	di.RegisterBean("IssueUseCases", reflect.TypeOf((*usecases.IssueUseCases)(nil)))

	di.RegisterBean("AccountController", reflect.TypeOf((*controllers.AccountController)(nil)))
	di.RegisterBean("ReportsController", reflect.TypeOf((*controllers.ReportsController)(nil)))
	di.RegisterBean("HealthController", reflect.TypeOf((*controllers.HealthController)(nil)))
	di.RegisterBean("ColumnController", reflect.TypeOf((*controllers.ColumnController)(nil)))
	di.RegisterBean("TeamController", reflect.TypeOf((*controllers.TeamController)(nil)))
	di.RegisterBean("LabelController", reflect.TypeOf((*controllers.LabelController)(nil)))
	di.RegisterBean("UserController", reflect.TypeOf((*controllers.UserController)(nil)))
	di.RegisterBean("GroupController", reflect.TypeOf((*controllers.GroupController)(nil)))
	di.RegisterBean("ProjectController", reflect.TypeOf((*controllers.ProjectController)(nil)))
	di.RegisterBean("SettingsController", reflect.TypeOf((*controllers.SettingsController)(nil)))
	di.RegisterBean("KanbanController", reflect.TypeOf((*controllers.KanbanController)(nil)))
	di.RegisterBean("SprintController", reflect.TypeOf((*controllers.SprintController)(nil)))
	di.RegisterBean("IssueController", reflect.TypeOf((*controllers.IssueController)(nil)))

	err = di.InitializeContainer()
	if err != nil {
		panic(err)
	}
}

func (app *Application) Run() {
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
