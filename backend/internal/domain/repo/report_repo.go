package repo

import domain "main/internal/domain/models"

type ReportRepo interface {
	GetBurndownData(sprintId uint) ([]domain.BurndownDataPoint, error)
}
