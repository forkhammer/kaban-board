package account_usecases

import (
	"errors"
	"main/internal/app/interfaces"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type RegisterUseCase struct {
	accountRepo     repo.AccountRepo                    `di.inject:"AccountRepository"`
	passwordService interfaces.PasswordServiceInterface `di.inject:"PasswordService"`
}

func (uc *RegisterUseCase) Execute(username, password string) (*domain.Account, error) {
	_, err := uc.accountRepo.GetByUsername(username)

	if err == nil {
		return nil, errors.New("Такой пользователь уже существует")
	}

	passwordHash, err := uc.passwordService.HashPassword(password)

	if err != nil {
		return nil, err
	}

	account := &domain.Account{
		Name:     username,
		Username: username,
		Password: passwordHash,
		IsActive: false,
	}
	account, err = uc.accountRepo.Create(account)

	return account, err
}
