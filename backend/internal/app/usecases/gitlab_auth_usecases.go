package usecases

import (
	"main/internal/app/interfaces"
	domain "main/internal/domain/models"
)

type GitLabAuthUseCases struct {
	gitLabAuthService interfaces.GitLabAuthServiceInterface `di.inject:"GitLabAuthService"`
	jwtService        interfaces.JWTServiceInterface        `di.inject:"JWTService"`
}

func (uc *GitLabAuthUseCases) GetAuthorizationURL() (string, error) {
	if !uc.gitLabAuthService.IsEnabled() {
		return "", nil
	}
	return uc.gitLabAuthService.GetAuthorizationURL(), nil
}

func (uc *GitLabAuthUseCases) Login(code string) (string, *domain.Account, error) {
	account, err := uc.gitLabAuthService.HandleCallback(code)
	if err != nil {
		return "", nil, err
	}

	token, err := uc.jwtService.GenerateToken(account)
	if err != nil {
		return "", nil, err
	}

	return token, account, nil
}

func (uc *GitLabAuthUseCases) IsEnabled() bool {
	return uc.gitLabAuthService.IsEnabled()
}
