package repo

import "main/internal/domain/models"

type GitlabTokenRepo interface {
	Upsert(token *models.GitlabToken) (*models.GitlabToken, error)
	GetByAccountId(accountId models.AccountId) (*models.GitlabToken, error)
}
