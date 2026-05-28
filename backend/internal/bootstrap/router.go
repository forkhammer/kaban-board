package bootstrap

import (
	"main/config"
	"main/internal/interfaces/api"
	"main/internal/interfaces/api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/goioc/di"
)

func InitRouter(router *gin.Engine) error {
	cfg := di.GetInstance("config").(*config.Config)

	router.Use(middleware.JwtMiddleware())
	router.Use(middleware.OnlineTrackingMiddleware(cfg))

	controllers := []api.Controller{
		di.GetInstance("AccountController").(api.Controller),
		di.GetInstance("GitLabAuthController").(api.Controller),
		di.GetInstance("KanbanController").(api.Controller),
		di.GetInstance("ReportsController").(api.Controller),
		di.GetInstance("HealthController").(api.Controller),
		di.GetInstance("ColumnController").(api.Controller),
		di.GetInstance("TeamController").(api.Controller),
		di.GetInstance("LabelController").(api.Controller),
		di.GetInstance("UserController").(api.Controller),
		di.GetInstance("GroupController").(api.Controller),
		di.GetInstance("ProjectController").(api.Controller),
		di.GetInstance("SettingsController").(api.Controller),
		di.GetInstance("SprintController").(api.Controller),
		di.GetInstance("IssueController").(api.Controller),
		di.GetInstance("IssueBindingController").(api.Controller),
		di.GetInstance("ReleaseController").(api.Controller),
		di.GetInstance("EpicController").(api.Controller),
		di.GetInstance("SprintWSController").(api.Controller),
	}

	apiGroup := router.Group("/api")
	for _, controller := range controllers {
		if err := controller.RegisterRoutes(apiGroup); err != nil {
			return err
		}
	}

	return nil
}
