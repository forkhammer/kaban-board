package interfaces

import (
	domain "main/internal/domain/models"
)

type JWTServiceInterface interface {
	GenerateToken(account *domain.Account) (string, error)
	ValidateToken(token string) error
	GetAccountId(token string) (domain.AccountId, error)
}
