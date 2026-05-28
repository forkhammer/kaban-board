package testutil

import (
	"context"
	"fmt"
	"main/config"
	"main/internal/bootstrap"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"
	"net/http/httptest"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/goioc/di"
)

type TestSuite struct {
	Server  *httptest.Server
	Config  *config.Config
	cleanup func()
}

func SetupSuiteForTestMain() *TestSuite {
	dbType := interfaces.DbType(os.Getenv("DB_TYPE"))
	if dbType == "" {
		dbType = interfaces.Sqlite
	}

	var cleanup func()
	ctx := context.Background()

	switch dbType {
	case interfaces.Postgresql:
		var err error
		cleanup, err = SetupPostgres(ctx)
		if err != nil {
			panic(fmt.Sprintf("failed to setup postgres: %v", err))
		}
	case interfaces.Mysql:
		var err error
		cleanup, err = SetupMysql(ctx)
		if err != nil {
			panic(fmt.Sprintf("failed to setup mysql: %v", err))
		}
	default:
		os.Setenv("SQLITE_DB_FILE", ":memory:")
	}

	os.Setenv("TEST_MODE", "true")
	cfg := config.GetConfig()

	bootstrap.InitDI()

	conn := di.GetInstance("connection").(interfaces.ConnectionInterface)
	if err := conn.Migrate(models.ALL_MODELS...); err != nil {
		if cleanup != nil {
			cleanup()
		}
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	router := setupTestRouter(cfg)
	server := httptest.NewServer(router)

	return &TestSuite{
		Server:  server,
		Config:  cfg,
		cleanup: cleanup,
	}
}

func (s *TestSuite) Teardown() {
	if s.Server != nil {
		s.Server.Close()
	}
	if s.cleanup != nil {
		s.cleanup()
	}
}

func setupTestRouter(cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())

	if err := bootstrap.InitRouter(router); err != nil {
		panic(fmt.Sprintf("failed to initialize router: %v", err))
	}

	return router
}
