package controllers

import (
	"main/internal/app/usecases"
	"main/internal/interfaces/api/dto"
	"main/internal/interfaces/api/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportsController struct {
	reportUC *usecases.ReportUseCases `di.inject:"ReportUseCases"`
}

func (c *ReportsController) RegisterRoutes(router gin.IRouter) error {
	router.GET("/reports/burndown", c.getBurndownReport)
	router.GET("/reports/burnup", c.getBurnupReport)
	router.GET("/reports/wip", c.getWipReport)
	router.GET("/reports/sprint/:id", c.getSprintStats)
	return nil
}

func (c *ReportsController) getSprintStats(ctx *gin.Context) {
	sprintId, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid sprint id"})
		return
	}

	var request dto.SprintStatsRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	stats, err := c.reportUC.GetSprintStats(uint(sprintId), request.AssigneeId, request.GroupId)
	if utils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeSprintStats(stats))
}

func (c *ReportsController) getWipReport(ctx *gin.Context) {
	var request dto.WipReportRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid start_date format, expected YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid end_date format, expected YYYY-MM-DD"})
		return
	}

	params := usecases.WipReportParams{
		StartDate: startDate,
		EndDate:   endDate,
		Interval:  request.Interval,
		TeamId:    request.TeamId,
		UserId:    request.UserId,
	}

	report, err := c.reportUC.GetWipReport(params)
	if utils.HandleException(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, dto.SerializeWipReport(report))
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
