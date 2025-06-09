package repo

import (
	"main/internal/domain/models"
)

type IssueBindingRepo interface {
	RWRepo[models.IssueBinding, models.IssueBindingId]
}
