package controllers

import (
	"main/internal/app/usecases"
	"main/internal/interfaces/api/dto"
	"main/internal/interfaces/api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportsController struct {
	reportUC *usecases.ReportUseCases `di.inject:"ReportUseCases"`
}

func (c *ReportsController) RegisterRoutes(router gin.IRouter) error {
	router.GET("/reports/burndown", c.getBurndownReport)
	router.GET("/reports/burnup", c.getBurnupReport)
	return nil
}

func (c *ReportsController) getBurndownReport(ctx *gin.Context) {
	var request dto.BurndownReportRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	report, err := c.reportUC.GetBurndownReport(request.SprintId)
	if utils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeBurndownReport(report))
}

func (c *ReportsController) getBurnupReport(ctx *gin.Context) {
	var request dto.BurnupReportRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	report, err := c.reportUC.GetBurnupReport(request.SprintId)
	if utils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeBurnupReport(report))
}
