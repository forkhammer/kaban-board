package repo

import (
	domain "main/internal/domain/models"
	"time"
)

type ReportRepo interface {
	GetBurndownData(sprintId uint) ([]domain.BurndownDataPoint, error)
	GetBurnupData(sprintId uint) ([]domain.BurnupDataPoint, error)
	GetWipData(startDate, endDate time.Time, teamId *uint, userId *uint) ([]domain.WipDataPoint, error)
}
