package usecases

import (
	"fmt"

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
