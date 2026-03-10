package repo

import (
	"main/internal/domain/models"
)

type AccountRepo interface {
	RWRepo[models.Account, models.AccountId]
	GetByUsername(username string) (*models.Account, error)
	GetByGitlabID(gitlabID uint) (*models.Account, error)
}
