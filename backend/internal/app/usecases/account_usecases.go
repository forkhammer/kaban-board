package usecases

import (
	"errors"
	"main/internal/app/interfaces"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type AccountUseCases struct {
	accountRepo     repo.AccountRepo                    `di.inject:"AccountRepository"`
	jwtService      interfaces.JWTServiceInterface      `di.inject:"JWTService"`
	passwordService interfaces.PasswordServiceInterface `di.inject:"PasswordService"`
}

func (uc *AccountUseCases) GetActiveUser(token string) (*domain.Account, error) {
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

func (uc *AccountUseCases) Login(username, password string) (string, error) {
	account, err := uc.accountRepo.GetByUsername(username)

	if err != nil {
		return "", errors.New("Такой пользователь не найден")
	}

	if !account.IsActive {
		return "", errors.New("Пользователь не активирован")
	}

	err = uc.passwordService.VerifyPassword(password, account.Password)

	if err != nil {
		return "", errors.New("Неверный пароль")
	}

	token, err := uc.jwtService.GenerateToken(account)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (uc *AccountUseCases) Register(username, password string) (*domain.Account, error) {
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
