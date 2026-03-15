package usecases

import (
	"main/internal/app/interfaces"
	domain_pkg "main/internal/domain"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
)

type AccountUseCases struct {
	accountRepo     repo.AccountRepo                    `di.inject:"AccountRepository"`
	jwtService      interfaces.JWTServiceInterface      `di.inject:"JWTService"`
	passwordService interfaces.PasswordServiceInterface `di.inject:"PasswordService"`
}

func (uc *AccountUseCases) GetActiveUser(token string) (*domain.Account, error) {
	accountId, err := uc.jwtService.GetAccountId(token)
	if err != nil {
		return nil, err
	}

	account, err := uc.accountRepo.Get(accountId)
	if err != nil {
		return nil, err
	}

	if err := uc.jwtService.ValidateToken(token, account.JwtSalt); err != nil {
		return nil, err
	}

	return account, nil
}

func (uc *AccountUseCases) Login(username, password string) (string, error) {
	account, err := uc.accountRepo.GetByUsername(username)

	if err != nil {
		return "", domain_pkg.NewValidationError("Пользователь не найден", nil)
	}

	if !account.IsActive {
		return "", domain_pkg.NewValidationError("Пользователь не активирон", nil)
	}

	err = uc.passwordService.VerifyPassword(password, account.Password)

	if err != nil {
		return "", domain_pkg.NewValidationError("Неверный пароль", nil)
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
		return nil, domain_pkg.NewValidationError("Такой пользователь не существует", nil)
	}

	passwordHash, err := uc.passwordService.HashPassword(password)
	if err != nil {
		return nil, err
	}

	salt, err := uc.passwordService.GenerateSalt()
	if err != nil {
		return nil, err
	}

	account := &domain.Account{
		Name:     username,
		Username: username,
		Password: passwordHash,
		IsActive: false,
		JwtSalt:  salt,
	}
	account, err = uc.accountRepo.Create(account)

	return account, err
}
