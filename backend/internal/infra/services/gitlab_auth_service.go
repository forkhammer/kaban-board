package services

import (
	"fmt"
	"main/config"
	"main/internal/app/interfaces"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/gitlab"
)

type GitLabAuthService struct {
	config       *config.Config       `di.inject:"config"`
	accountRepo  repo.AccountRepo     `di.inject:"AccountRepository"`
	gitlabClient *gitlab.GitlabClient `di.inject:"gitlab"`
}

func (s *GitLabAuthService) IsEnabled() bool {
	return s.config.GitLabAuthEnabled
}

func (s *GitLabAuthService) GetAuthorizationURL() string {
	authURL := fmt.Sprintf("%s/oauth/authorize?client_id=%s&response_type=code&scope=%s&state=gitlab_auth&redirect_uri=%s",
		s.config.GitlabUrl,
		s.config.GitLabAuthClientID,
		s.config.GitLabAuthScope,
		s.config.GitLabAuthRedirectURL,
	)
	return authURL
}

func (s *GitLabAuthService) HandleCallback(code string) (*domain.Account, error) {
	token, err := s.gitlabClient.ExchangeCodeForToken(code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	userInfo, err := s.gitlabClient.GetUserInfo(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	account, err := s.findOrCreateAccount(userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create account: %w", err)
	}

	return account, nil
}

func (s *GitLabAuthService) findOrCreateAccount(userInfo *interfaces.GitLabUserInfo) (*domain.Account, error) {
	account, err := s.accountRepo.GetByGitlabID(userInfo.ID)
	if err == nil && account != nil {
		account.AvatarURL = userInfo.AvatarURL
		return s.accountRepo.Update(account)
	}

	gitlabId := domain.UserId(userInfo.ID)
	account = &domain.Account{
		Username:     userInfo.Username,
		Name:         userInfo.Name,
		GitlabID:     &gitlabId,
		AuthProvider: domain.AuthProviderGitLab,
		AvatarURL:    userInfo.AvatarURL,
		IsActive:     true,
	}

	return s.accountRepo.Create(account)
}
