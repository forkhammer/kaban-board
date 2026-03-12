package dto

import domain "main/internal/domain/models"

type BurndownReportRequest struct {
	SprintId uint `form:"sprint_id" binding:"required"`
}

type BurndownDataPointDto struct {
	Date         string `json:"date"`
	TotalScope   uint   `json:"total_scope"`
	RemainingDev uint   `json:"remaining_dev"`
}

type BurndownReportDto struct {
	SprintTitle string                 `json:"sprint_title"`
	StartDate   string                 `json:"start_date"`
	EndDate     string                 `json:"end_date"`
	DataPoints  []BurndownDataPointDto `json:"data_points"`
}

func SerializeBurndownReport(report *domain.BurndownReport) BurndownReportDto {
	dataPoints := make([]BurndownDataPointDto, len(report.DataPoints))
	for i, dp := range report.DataPoints {
		dataPoints[i] = BurndownDataPointDto{
			Date:         dp.Date.Format("02.01.2006"),
			TotalScope:   dp.TotalScope,
			RemainingDev: dp.RemainingDev,
		}
	}

	return BurndownReportDto{
		SprintTitle: report.SprintTitle,
		StartDate:   report.StartDate.Format("2006-01-02"),
		EndDate:     report.EndDate.Format("2006-01-02"),
		DataPoints:  dataPoints,
	}
}
