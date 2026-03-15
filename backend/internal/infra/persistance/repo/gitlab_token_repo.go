package repo

import (
	"errors"
	"log"
	"strconv"

	app_interfaces "main/internal/app/interfaces"
	domain_pkg "main/internal/domain"
	domain "main/internal/domain/models"
	"main/internal/infra/db/interfaces"
	"main/internal/infra/persistance/models"

	"github.com/getsentry/sentry-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GitlabTokenRepository struct {
	conn          interfaces.ConnectionInterface        `di.inject:"db"`
	cryptoService app_interfaces.CryptoServiceInterface `di.inject:"CryptoService"`
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
	accessToken, err := r.cryptoService.Decrypt(m.AccessToken)
	if err != nil {
		log.Printf("failed to decrypt access_token for account %d: %v", m.AccountId, err)
		accessToken = m.AccessToken
	}
	refreshToken, err := r.cryptoService.Decrypt(m.RefreshToken)
	if err != nil {
		log.Printf("failed to decrypt refresh_token for account %d: %v", m.AccountId, err)
		refreshToken = m.RefreshToken
	}
	return &domain.GitlabToken{
		Id:           domain.GitlabTokenId(m.Id),
		AccountId:    domain.AccountId(m.AccountId),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    m.ExpiresAt,
	}
}

func (r *GitlabTokenRepository) toModel(t *domain.GitlabToken) *models.GitlabToken {
	accessToken, err := r.cryptoService.Encrypt(t.AccessToken)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("failed to encrypt access_token for account %d: %v", t.AccountId, err)
		accessToken = t.AccessToken
	}
	refreshToken, err := r.cryptoService.Encrypt(t.RefreshToken)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("failed to encrypt refresh_token for account %d: %v", t.AccountId, err)
		refreshToken = t.RefreshToken
	}
	return &models.GitlabToken{
		Id:           uint(t.Id),
		AccountId:    uint(t.AccountId),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    t.ExpiresAt,
	}
}
