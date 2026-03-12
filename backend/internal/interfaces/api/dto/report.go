package dto

import domain "main/internal/domain/models"

type BurndownReportRequest struct {
	SprintId uint `form:"sprint_id" binding:"required"`
}

type BurndownDataPointDto struct {
	Date         string `json:"date"`
	TotalScope   uint   `json:"total_scope"`
	RemainingDev uint   `json:"remaining_dev"`
	TotalScopeQA uint   `json:"total_scope_qa"`
	RemainingQA  uint   `json:"remaining_qa"`
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
			TotalScopeQA: dp.TotalScopeQA,
			RemainingQA:  dp.RemainingQA,
		}
	}

	return BurndownReportDto{
		SprintTitle: report.SprintTitle,
		StartDate:   report.StartDate.Format("2006-01-02"),
		EndDate:     report.EndDate.Format("2006-01-02"),
		DataPoints:  dataPoints,
	}
}

// Burnup Report DTOs

type BurnupReportRequest struct {
	SprintId uint `form:"sprint_id" binding:"required"`
}

type BurnupDataPointDto struct {
	Date         string `json:"date"`
	ScopeDev     uint   `json:"scope_dev"`
	CompletedDev uint   `json:"completed_dev"`
	ScopeQA      uint   `json:"scope_qa"`
	CompletedQA  uint   `json:"completed_qa"`
}

type BurnupReportDto struct {
	SprintTitle string               `json:"sprint_title"`
	StartDate   string               `json:"start_date"`
	EndDate     string               `json:"end_date"`
	DataPoints  []BurnupDataPointDto `json:"data_points"`
}

func SerializeBurnupReport(report *domain.BurnupReport) BurnupReportDto {
	dataPoints := make([]BurnupDataPointDto, len(report.DataPoints))
	for i, dp := range report.DataPoints {
		dataPoints[i] = BurnupDataPointDto{
			Date:         dp.Date.Format("02.01.2006"),
			ScopeDev:     dp.ScopeDev,
			CompletedDev: dp.CompletedDev,
			ScopeQA:      dp.ScopeQA,
			CompletedQA:  dp.CompletedQA,
		}
	}

	return BurnupReportDto{
		SprintTitle: report.SprintTitle,
		StartDate:   report.StartDate.Format("2006-01-02"),
		EndDate:     report.EndDate.Format("2006-01-02"),
		DataPoints:  dataPoints,
	}
}
