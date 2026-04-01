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
	Interval  string
	TeamId    *uint
	UserId    *uint
}

var ruMonths = [12]string{"Янв", "Фев", "Мар", "Апр", "Май", "Июн", "Июл", "Авг", "Сен", "Окт", "Ноя", "Дек"}

func wipLabel(date time.Time, interval string) string {
	switch interval {
	case "week":
		_, week := date.ISOWeek()
		return fmt.Sprintf("Нед. %d", week)
	case "2weeks":
		_, week := date.ISOWeek()
		return fmt.Sprintf("%d-%d", week, week+1)
	case "month":
		return fmt.Sprintf("%s %d", ruMonths[date.Month()-1], date.Year())
	default: // day
		return date.Format("02.01.2006")
	}
}

func (uc *ReportUseCases) GetWipReport(params WipReportParams) (*domain.WipReport, error) {
	dataPoints, err := uc.reportRepo.GetWipData(params.StartDate, params.EndDate, params.Interval, params.TeamId, params.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to get wip data: %w", err)
	}

	for i := range dataPoints {
		dataPoints[i].Label = wipLabel(dataPoints[i].Date, params.Interval)
	}

	return &domain.WipReport{DataPoints: dataPoints}, nil
}

func (uc *ReportUseCases) GetSprintStats(sprintId uint, assigneeId *uint, groupId *uint) (*domain.SprintStats, error) {
	sprint, err := uc.sprintRepo.Get(domain.SprintId(sprintId))
	if err != nil {
		return nil, fmt.Errorf("sprint not found: %w", err)
	}

	var userCount uint
	if assigneeId != nil {
		userCount = 1
	} else {
		userCount, err = uc.reportRepo.GetActiveUserCountByTeam(uint(sprint.Team.Id), groupId)
		if err != nil {
			return nil, fmt.Errorf("failed to count team members: %w", err)
		}
	}

	stats, err := uc.reportRepo.GetSprintStats(sprintId, assigneeId, groupId)
	if err != nil {
		return nil, fmt.Errorf("failed to get sprint stats: %w", err)
	}

	stats.Capacity = sprint.HoursPerUser * userCount
	return stats, nil
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
