package repo

import (
	"errors"
	"strconv"

	domain_pkg "main/internal/domain"
	domain "main/internal/domain/models"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GitlabTokenRepository struct {
	conn interfaces.ConnectionInterface `di.inject:"db"`
}

func (r *GitlabTokenRepository) GetByAccountId(accountId domain.AccountId) (*domain.GitlabToken, error) {
	model := &models.GitlabToken{}
	if err := r.conn.GetEngine().Where("account_id = ?", accountId).First(model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain_pkg.NewNotFoundError("GitlabToken", strconv.Itoa(int(accountId)), err)
		}
		return nil, err
	}
	return r.toDomain(model), nil
}

func (r *GitlabTokenRepository) Upsert(token *domain.GitlabToken) (*domain.GitlabToken, error) {
	model := r.toModel(token)
	err := r.conn.GetEngine().
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "account_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"access_token", "refresh_token", "expires_at", "updated_at"}),
		}).
		Create(model).Error
	if err != nil {
		return nil, err
	}
	return r.GetByAccountId(domain.AccountId(model.AccountId))
}

func (r *GitlabTokenRepository) toDomain(m *models.GitlabToken) *domain.GitlabToken {
	return &domain.GitlabToken{
		Id:           domain.GitlabTokenId(m.Id),
		AccountId:    domain.AccountId(m.AccountId),
		AccessToken:  m.AccessToken,
		RefreshToken: m.RefreshToken,
		ExpiresAt:    m.ExpiresAt,
	}
}

func (r *GitlabTokenRepository) toModel(t *domain.GitlabToken) *models.GitlabToken {
	return &models.GitlabToken{
		Id:           uint(t.Id),
		AccountId:    uint(t.AccountId),
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		ExpiresAt:    t.ExpiresAt,
	}
}
