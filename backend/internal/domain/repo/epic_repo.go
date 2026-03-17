package repo

import (
	"main/internal/domain/models"
)

type EpicRepo interface {
	RWRepo[models.Epic, models.EpicId]
	GetByExternalId(externalId models.IssueExternalId) (*models.Epic, error)
}
