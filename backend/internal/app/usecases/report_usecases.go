package usecases

import (
	"fmt"
	"time"

	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type ReportUseCases struct {
	reportRepo repo.ReportRepo `di.inject:"ReportRepository"`
	sprintRepo repo.SprintRepo `di.inject:"SprintRepository"`
}

func (uc *ReportUseCases) GetBurndownReport(sprintId uint) (*domain.BurndownReport, error) {
	sprint, err := uc.sprintRepo.Get(domain.SprintId(sprintId))
	if err != nil {
		return nil, fmt.Errorf("sprint not found: %w", err)
	}

	dataPoints, err := uc.reportRepo.GetBurndownData(sprintId)
	if err != nil {
		return nil, fmt.Errorf("failed to get burndown data: %w", err)
	}

	report := &domain.BurndownReport{
		SprintTitle: sprint.GetTitle(),
		StartDate:   sprint.StartDate,
		EndDate:     sprint.EndDate,
		DataPoints:  dataPoints,
	}

	return report, nil
}

type WipReportParams struct {
	StartDate time.Time
	EndDate   time.Time
	TeamId    *uint
	UserId    *uint
}

func (uc *ReportUseCases) GetWipReport(params WipReportParams) (*domain.WipReport, error) {
	dataPoints, err := uc.reportRepo.GetWipData(params.StartDate, params.EndDate, params.TeamId, params.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to get wip data: %w", err)
	}

	return &domain.WipReport{DataPoints: dataPoints}, nil
}

func (uc *ReportUseCases) GetBurnupReport(sprintId uint) (*domain.BurnupReport, error) {
	sprint, err := uc.sprintRepo.Get(domain.SprintId(sprintId))
	if err != nil {
		return nil, fmt.Errorf("sprint not found: %w", err)
	}

	dataPoints, err := uc.reportRepo.GetBurnupData(sprintId)
	if err != nil {
		return nil, fmt.Errorf("failed to get burnup data: %w", err)
	}

	report := &domain.BurnupReport{
		SprintTitle: sprint.GetTitle(),
		StartDate:   sprint.StartDate,
		EndDate:     sprint.EndDate,
		DataPoints:  dataPoints,
	}

	return report, nil
}
