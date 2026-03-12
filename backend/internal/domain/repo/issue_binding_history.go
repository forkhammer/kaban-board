package repo

import (
	"main/internal/domain/models"
)

type IssueBindingHistoryRepo interface {
	RetrieveRepo[models.IssueBindingHistory, models.IssueBindingHistoryId]
	ListRepo[models.IssueBindingHistory]
	CreateRepo[models.IssueBindingHistory]
	DeleteRepo[models.IssueBindingHistoryId]
}
