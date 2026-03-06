package main

import (
	"context"
	"fmt"
	"main/cmd"
	"main/config"
	app_services "main/internal/app/services"
	"main/internal/app/usecases"
	"main/internal/infra/db/implementation"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/gitlab"
	"main/internal/infra/persistance/models"
	"main/internal/infra/persistance/repo"
	"main/internal/infra/persistance/spec"
	column_spec "main/internal/infra/persistance/spec/column"
	epic_spec "main/internal/infra/persistance/spec/epic"
	group_spec "main/internal/infra/persistance/spec/group"
	issue_spec "main/internal/infra/persistance/spec/issue"
	issuebinding_spec "main/internal/infra/persistance/spec/issue_binding"
	label_spec "main/internal/infra/persistance/spec/label"
	project_spec "main/internal/infra/persistance/spec/project"
	release_spec "main/internal/infra/persistance/spec/release"
	sprint_spec "main/internal/infra/persistance/spec/sprint"
	user_spec "main/internal/infra/persistance/spec/user"
	"main/internal/infra/services"
	"main/internal/interfaces/api/controllers"
	"reflect"
	"time"

	sentry "github.com/getsentry/sentry-go"
	"github.com/goioc/di"
)

func mustRegisterBean(beanID string, beanType reflect.Type) {
	if _, err := di.RegisterBean(beanID, beanType); err != nil {
		panic(fmt.Sprintf("failed to register bean %q: %v", beanID, err))
	}
}

func mustRegisterBeanInstance(beanID string, beanInstance any) {
	if _, err := di.RegisterBeanInstance(beanID, beanInstance); err != nil {
		panic(fmt.Sprintf("failed to register bean instance %q: %v", beanID, err))
	}
}

func mustRegisterBeanFactory(beanID string, scope di.Scope, beanFactory func(ctx context.Context) (any, error)) {
	if _, err := di.RegisterBeanFactory(beanID, scope, beanFactory); err != nil {
		panic(fmt.Sprintf("failed to register bean factory %q: %v", beanID, err))
	}
}

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

	connection, err := implementation.GetConnectionByType(interfaces.DbType(config.Settings.DbType), config.Settings)
	if err != nil {
		panic(err)
	}

	mustRegisterBeanInstance("connection", connection)
	mustRegisterBeanInstance("config", config.Settings)
	mustRegisterBeanInstance("db", connection)
	mustRegisterBeanFactory("gitlab", di.Singleton, func(ctx context.Context) (any, error) {
		return gitlab.NewGitlabClient(config.Settings.GitlabUrl, config.Settings.GitlabToken), nil
	})

	mustRegisterBean("AccountRepository", reflect.TypeFor[*repo.AccountRepository]())
	mustRegisterBean("ColumnRepository", reflect.TypeFor[*repo.ColumnRepository]())
	mustRegisterBean("GroupRepository", reflect.TypeFor[*repo.GroupRepository]())
	mustRegisterBean("IssueRepository", reflect.TypeFor[*repo.IssueRepository]())
	mustRegisterBean("KeyValueRepository", reflect.TypeFor[*repo.KeyValueRepository]())
	mustRegisterBean("LabelRepository", reflect.TypeFor[*repo.LabelRepository]())
	mustRegisterBean("ProjectRepository", reflect.TypeFor[*repo.ProjectRepository]())
	mustRegisterBean("ReleaseRepository", reflect.TypeFor[*repo.ReleaseRepository]())
	mustRegisterBean("TeamRepository", reflect.TypeFor[*repo.TeamRepository]())
	mustRegisterBean("UserRepository", reflect.TypeFor[*repo.UserRepository]())
	mustRegisterBean("SettingsRepository", reflect.TypeFor[*repo.SettingsRepository]())
	mustRegisterBean("SprintRepository", reflect.TypeFor[*repo.SprintRepository]())
	mustRegisterBean("IssueBindingRepository", reflect.TypeFor[*repo.IssueBindingRepository]())
	mustRegisterBean("EpicRepository", reflect.TypeFor[*repo.EpicRepository]())

	mustRegisterBean("LabelQuery", reflect.TypeFor[*label_spec.LabelQueryImpl]())
	mustRegisterBean("ColumnQuery", reflect.TypeFor[*column_spec.ColumnQueryImpl]())
	mustRegisterBean("GroupQuery", reflect.TypeFor[*group_spec.GroupQueryImpl]())
	mustRegisterBean("UserQuery", reflect.TypeFor[*user_spec.UserQueryImpl]())
	mustRegisterBean("SprintQuery", reflect.TypeFor[*sprint_spec.SprintQueryImpl]())
	mustRegisterBean("IssueQuery", reflect.TypeFor[*issue_spec.IssueQueryImpl]())
	mustRegisterBean("ProjectQuery", reflect.TypeFor[*project_spec.ProjectQueryImpl]())
	mustRegisterBean("IssueBindingQuery", reflect.TypeFor[*issuebinding_spec.IssueBindingQueryImpl]())
	mustRegisterBean("ReleaseQuery", reflect.TypeFor[*release_spec.ReleaseQueryImpl]())
	mustRegisterBean("EpicQuery", reflect.TypeFor[*epic_spec.EpicQueryImpl]())
	mustRegisterBean("CommonQuery", reflect.TypeFor[*spec.CommonQueryImpl]())

	mustRegisterBean("JWTService", reflect.TypeFor[*services.JWTService]())
	mustRegisterBean("PasswordService", reflect.TypeFor[*services.PasswordService]())
	mustRegisterBean("LabelService", reflect.TypeFor[*app_services.LabelService]())

	mustRegisterBean("AccountUseCases", reflect.TypeFor[*usecases.AccountUseCases]())
	mustRegisterBean("ColumnUseCases", reflect.TypeFor[*usecases.ColumnUseCases]())
	mustRegisterBean("TeamUseCases", reflect.TypeFor[*usecases.TeamUseCases]())
	mustRegisterBean("LabelUseCases", reflect.TypeFor[*usecases.LabelUseCases]())
	mustRegisterBean("UserUseCases", reflect.TypeFor[*usecases.UserUseCases]())
	mustRegisterBean("GroupUseCases", reflect.TypeFor[*usecases.GroupUseCases]())
	mustRegisterBean("ProjectUseCases", reflect.TypeFor[*usecases.ProjectUseCases]())
	mustRegisterBean("SettingsUseCases", reflect.TypeFor[*usecases.SettingsUseCases]())
	mustRegisterBean("SyncUseCases", reflect.TypeFor[*usecases.SyncUseCases]())
	mustRegisterBean("KanbanUseCases", reflect.TypeFor[*usecases.KanbanUseCases]())
	mustRegisterBean("SprintUseCases", reflect.TypeFor[*usecases.SprintUseCases]())
	mustRegisterBean("IssueUseCases", reflect.TypeFor[*usecases.IssueUseCases]())
	mustRegisterBean("IssueBindingUseCases", reflect.TypeFor[*usecases.IssueBindingUseCases]())
	mustRegisterBean("ReleaseUseCases", reflect.TypeFor[*usecases.ReleaseUseCases]())
	mustRegisterBean("EpicUseCases", reflect.TypeFor[*usecases.EpicUseCases]())

	mustRegisterBean("AccountController", reflect.TypeFor[*controllers.AccountController]())
	mustRegisterBean("ReportsController", reflect.TypeFor[*controllers.ReportsController]())
	mustRegisterBean("HealthController", reflect.TypeFor[*controllers.HealthController]())
	mustRegisterBean("ColumnController", reflect.TypeFor[*controllers.ColumnController]())
	mustRegisterBean("TeamController", reflect.TypeFor[*controllers.TeamController]())
	mustRegisterBean("LabelController", reflect.TypeFor[*controllers.LabelController]())
	mustRegisterBean("UserController", reflect.TypeFor[*controllers.UserController]())
	mustRegisterBean("GroupController", reflect.TypeFor[*controllers.GroupController]())
	mustRegisterBean("ProjectController", reflect.TypeFor[*controllers.ProjectController]())
	mustRegisterBean("SettingsController", reflect.TypeFor[*controllers.SettingsController]())
	mustRegisterBean("KanbanController", reflect.TypeFor[*controllers.KanbanController]())
	mustRegisterBean("SprintController", reflect.TypeFor[*controllers.SprintController]())
	mustRegisterBean("IssueController", reflect.TypeFor[*controllers.IssueController]())
	mustRegisterBean("IssueBindingController", reflect.TypeFor[*controllers.IssueBindingController]())
	mustRegisterBean("ReleaseController", reflect.TypeFor[*controllers.ReleaseController]())
	mustRegisterBean("EpicController", reflect.TypeFor[*controllers.EpicController]())

	err = di.InitializeContainer()
	if err != nil {
		panic(err)
	}
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
