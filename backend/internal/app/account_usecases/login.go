package account_usecases

import (
	"errors"
	"main/internal/app/interfaces"
	"main/internal/domain/repo"
)

type LoginUseCase struct {
	jwtService      interfaces.JWTServiceInterface      `di.inject:"JWTService"`
	accountRepo     repo.AccountRepo                    `di.inject:"AccountRepository"`
	passwordService interfaces.PasswordServiceInterface `di.inject:"PasswordService"`
}

func (uc *LoginUseCase) Execute(username, password string) (string, error) {
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
