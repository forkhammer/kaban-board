package repo

import (
	domain "main/internal/domain/models"
	"time"
)

type ReportRepo interface {
	GetBurndownData(sprintId uint) ([]domain.BurndownDataPoint, error)
	GetBurnupData(sprintId uint) ([]domain.BurnupDataPoint, error)
	GetWipData(startDate, endDate time.Time, interval string, teamId *uint, userId *uint) ([]domain.WipDataPoint, error)
}
