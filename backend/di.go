package main

import (
	"context"
	"fmt"
	"main/config"
	app_services "main/internal/app/services"
	"main/internal/app/usecases"
	"main/internal/infra/cache"
	"main/internal/infra/db/implementation"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/gitlab"
	"main/internal/infra/hub"
	"main/internal/infra/persistance/repo"
	"main/internal/infra/persistance/spec"
	account_spec "main/internal/infra/persistance/spec/account"
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

func initDI() {
	dbType := interfaces.DbType(config.Settings.DbType)

	connection, err := implementation.GetConnectionByType(dbType, config.Settings)
	if err != nil {
		panic(err)
	}

	mustRegisterBeanInstance("connection", connection)
	mustRegisterBeanInstance("config", config.Settings)
	mustRegisterBeanInstance("db", connection)
	mustRegisterBeanInstance("cache", cache.MemoryCacheInstance)
	mustRegisterBeanFactory("gitlab", di.Singleton, func(ctx context.Context) (any, error) {
		return gitlab.NewGitlabClient(config.Settings.GitlabUrl, config.Settings.GitlabToken), nil
	})

	mustRegisterBeanInstance("SprintHub", hub.NewSprintHub())

	registerRepositories()
	registerReportRepository(dbType)
	registerQueries(dbType)
	registerServices()
	registerUseCases()
	registerControllers()

	if err = di.InitializeContainer(); err != nil {
		panic(err)
	}
}

func registerRepositories() {
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
	mustRegisterBean("IssueBindingHistoryRepository", reflect.TypeFor[*repo.IssueBindingHistoryRepository]())
	mustRegisterBean("EpicRepository", reflect.TypeFor[*repo.EpicRepository]())
	mustRegisterBean("GitlabTokenRepository", reflect.TypeFor[*repo.GitlabTokenRepository]())
}

func registerReportRepository(dbType interfaces.DbType) {
	switch dbType {
	case interfaces.Postgresql:
		mustRegisterBean("ReportRepository", reflect.TypeFor[*repo.ReportRepositoryPostgresql]())
	default:
		mustRegisterBean("ReportRepository", reflect.TypeFor[*repo.ReportRepository]())
	}
}

func registerQueries(dbType interfaces.DbType) {
	mustRegisterBean("LabelQuery", reflect.TypeFor[*label_spec.LabelQueryImpl]())
	mustRegisterBean("ColumnQuery", reflect.TypeFor[*column_spec.ColumnQueryImpl]())
	mustRegisterBean("GroupQuery", reflect.TypeFor[*group_spec.GroupQueryImpl]())
	registerUserQuery(dbType)
	mustRegisterBean("AccountQuery", reflect.TypeFor[*account_spec.AccountQueryImpl]())
	mustRegisterBean("SprintQuery", reflect.TypeFor[*sprint_spec.SprintQueryImpl]())
	registerProjectQuery(dbType)
	mustRegisterBean("IssueBindingQuery", reflect.TypeFor[*issuebinding_spec.IssueBindingQueryImpl]())
	mustRegisterBean("ReleaseQuery", reflect.TypeFor[*release_spec.ReleaseQueryImpl]())
	registerEpicQuery(dbType)
	mustRegisterBean("CommonQuery", reflect.TypeFor[*spec.CommonQueryImpl]())

	registerIssueQuery(dbType)
}

func registerEpicQuery(dbType interfaces.DbType) {
	switch dbType {
	case interfaces.Postgresql:
		mustRegisterBean("EpicQuery", reflect.TypeFor[*epic_spec.EpicQueryImplPostgresql]())
	case interfaces.Mysql:
		mustRegisterBean("EpicQuery", reflect.TypeFor[*epic_spec.EpicQueryImplMysql]())
	default:
		mustRegisterBean("EpicQuery", reflect.TypeFor[*epic_spec.EpicQueryImpl]())
	}
}

func registerIssueQuery(dbType interfaces.DbType) {
	switch dbType {
	case interfaces.Postgresql:
		mustRegisterBean("IssueQuery", reflect.TypeFor[*issue_spec.IssueQueryImplPostgresql]())
	case interfaces.Mysql:
		mustRegisterBean("IssueQuery", reflect.TypeFor[*issue_spec.IssueQueryImplMysql]())
	default:
		mustRegisterBean("IssueQuery", reflect.TypeFor[*issue_spec.IssueQueryImpl]())
	}
}

func registerUserQuery(dbType interfaces.DbType) {
	switch dbType {
	case interfaces.Postgresql:
		mustRegisterBean("UserQuery", reflect.TypeFor[*user_spec.UserQueryImplPostgresql]())
	case interfaces.Mysql:
		mustRegisterBean("UserQuery", reflect.TypeFor[*user_spec.UserQueryImplMysql]())
	default:
		mustRegisterBean("UserQuery", reflect.TypeFor[*user_spec.UserQueryImpl]())
	}
}

func registerProjectQuery(dbType interfaces.DbType) {
	switch dbType {
	case interfaces.Postgresql:
		mustRegisterBean("ProjectQuery", reflect.TypeFor[*project_spec.ProjectQueryImplPostgresql]())
	case interfaces.Mysql:
		mustRegisterBean("ProjectQuery", reflect.TypeFor[*project_spec.ProjectQueryImplMysql]())
	default:
		mustRegisterBean("ProjectQuery", reflect.TypeFor[*project_spec.ProjectQueryImpl]())
	}
}

func registerServices() {
	mustRegisterBean("CryptoService", reflect.TypeFor[*services.CryptoService]())
	mustRegisterBean("JWTService", reflect.TypeFor[*services.JWTService]())
	mustRegisterBean("PasswordService", reflect.TypeFor[*services.PasswordService]())
	mustRegisterBean("LabelService", reflect.TypeFor[*app_services.LabelService]())
	mustRegisterBean("IssueBindingHistoryService", reflect.TypeFor[*app_services.IssueBindingHistoryService]())
	mustRegisterBean("GitLabAuthService", reflect.TypeFor[*services.GitLabAuthService]())
}

func registerUseCases() {
	mustRegisterBean("AccountUseCases", reflect.TypeFor[*usecases.AccountUseCases]())
	mustRegisterBean("GitLabAuthUseCases", reflect.TypeFor[*usecases.GitLabAuthUseCases]())
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
	mustRegisterBean("ReportUseCases", reflect.TypeFor[*usecases.ReportUseCases]())
}

func registerControllers() {
	mustRegisterBean("AccountController", reflect.TypeFor[*controllers.AccountController]())
	mustRegisterBean("GitLabAuthController", reflect.TypeFor[*controllers.GitLabAuthController]())
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
	mustRegisterBean("SprintWSController", reflect.TypeFor[*controllers.SprintWSController]())
}
