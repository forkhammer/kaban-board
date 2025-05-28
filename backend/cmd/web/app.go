package web

import (
	"main/config"
	"main/internal/app/account_usecases"
	"main/internal/app/column_usecases"
	"main/internal/infra/db/implementation"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"
	"main/internal/infra/persistance/repo"
	"main/internal/infra/services"
	"main/internal/interfaces/api"
	"main/internal/interfaces/api/controllers"
	"main/internal/interfaces/api/middleware"
	"reflect"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/goioc/di"
)

type Application struct {
	router     *gin.Engine
	connection interfaces.ConnectionInterface
}

func NewApplication() *Application {
	config.Settings.Print()

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.Settings.AllowOrigins
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	router.Use(cors.New(corsConfig))
	router.Use(middleware.JwtMiddleware())

	connection, err := implementation.GetConnectionByType(interfaces.DbType(config.Settings.DbType), config.Settings)
	di.RegisterBeanInstance("connection", connection)

	if err != nil {
		panic(err)
	}

	app := Application{
		router:     router,
		connection: connection,
	}
	app.registerDeps()
	app.migrate()
	return &app
}

func (a *Application) registerDeps() {
	di.RegisterBeanInstance("config", config.Settings)
	di.RegisterBeanInstance("db", a.connection)
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

	di.RegisterBean("JWTService", reflect.TypeOf((*services.JWTService)(nil)))
	di.RegisterBean("PasswordService", reflect.TypeOf((*services.PasswordService)(nil)))

	di.RegisterBean("ActiveUserUseCase", reflect.TypeOf((*account_usecases.ActiveUserUseCase)(nil)))
	di.RegisterBean("LoginUseCase", reflect.TypeOf((*account_usecases.LoginUseCase)(nil)))
	di.RegisterBean("RegisterUseCase", reflect.TypeOf((*account_usecases.RegisterUseCase)(nil)))

	di.RegisterBean("ListColumnsUseCase", reflect.TypeOf((*column_usecases.ListColumnsUseCase)(nil)))

	di.RegisterBean("AccountController", reflect.TypeOf((*controllers.AccountController)(nil)))
	di.RegisterBean("BoardController", reflect.TypeOf((*controllers.BoardController)(nil)))
	di.RegisterBean("ReportsController", reflect.TypeOf((*controllers.ReportsController)(nil)))
	di.RegisterBean("HealthController", reflect.TypeOf((*controllers.HealthController)(nil)))

	di.RegisterBean("ColumnController", reflect.TypeOf((*controllers.ColumnController)(nil)))
}

func (app *Application) Run() {
	di.InitializeContainer()

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
	}
	for _, controller := range controllers {
		err := controller.RegisterRoutes(app.router)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *Application) migrate() error {
	return app.connection.Migrate(models.ALL_MODELS...)
}
