package account_usecases

import (
	"main/internal/app/interfaces"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type ActiveUserUseCase struct {
	accountRepo repo.AccountRepo               `di.inject:"AccountRepository"`
	jwtService  interfaces.JWTServiceInterface `di.inject:"JWTService"`
}

func (uc *ActiveUserUseCase) Execute(token string) (*domain.Account, error) {
	err := uc.jwtService.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	accountId, err := uc.jwtService.GetAccountId(token)
	if err != nil {
		return nil, err
	}

	account, err := uc.accountRepo.Get(accountId)
	if err != nil {
		return nil, err
	}
	return account, nil
}
