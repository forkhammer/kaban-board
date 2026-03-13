package repo

import (
	"main/internal/domain/models"
)

type IssueRepo interface {
	RWRepo[models.Issue, models.IssueId]
	GetByExternalId(externalId models.IssueExternalId) (*models.Issue, error)
}
