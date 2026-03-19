package services

import (
	"fmt"
	"log"
	"main/config"
	"main/internal/app/interfaces"
	domain "main/internal/domain/models"
	"main/internal/domain/repo"
	"main/internal/infra/gitlab"
	"time"

	"github.com/getsentry/sentry-go"
)

type GitLabAuthService struct {
	config          *config.Config                    `di.inject:"config"`
	accountRepo     repo.AccountRepo                  `di.inject:"AccountRepository"`
	gitlabTokenRepo repo.GitlabTokenRepo              `di.inject:"GitlabTokenRepository"`
	passwordService interfaces.PasswordServiceInterface `di.inject:"PasswordService"`
	gitlabClient    *gitlab.GitlabClient              `di.inject:"gitlab"`
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

	expiresAt := time.Unix(int64(token.CreatedAt), 0).Add(time.Duration(token.ExpiresIn) * time.Second)
	if _, err := s.gitlabTokenRepo.Upsert(&domain.GitlabToken{
		AccountId:    account.Id,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    expiresAt,
	}); err != nil {
		sentry.CaptureException(err)
		log.Printf("failed to save gitlab token for account %d: %v", account.Id, err)
	}

	return account, nil
}

func (s *GitLabAuthService) GetValidToken(accountId domain.AccountId) (*domain.GitlabToken, error) {
	token, err := s.gitlabTokenRepo.GetByAccountId(accountId)
	if err != nil {
		return nil, err
	}

	if token.ExpiresAt.Before(time.Now().Add(5 * time.Minute)) {
		newToken, err := s.gitlabClient.RefreshToken(token.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh GitLab token: %w", err)
		}
		expiresAt := time.Unix(int64(newToken.CreatedAt), 0).Add(time.Duration(newToken.ExpiresIn) * time.Second)
		token, err = s.gitlabTokenRepo.Upsert(&domain.GitlabToken{
			AccountId:    accountId,
			AccessToken:  newToken.AccessToken,
			RefreshToken: newToken.RefreshToken,
			ExpiresAt:    expiresAt,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to save refreshed token: %w", err)
		}
	}

	return token, nil
}

func (s *GitLabAuthService) findOrCreateAccount(userInfo *interfaces.GitLabUserInfo) (*domain.Account, error) {
	account, err := s.accountRepo.GetByGitlabID(userInfo.ID)
	if err == nil && account != nil {
		account.AvatarURL = userInfo.AvatarURL
		return s.accountRepo.Update(account)
	}

	salt, err := s.passwordService.GenerateSalt()
	if err != nil {
		return nil, fmt.Errorf("failed to generate jwt salt: %w", err)
	}

	gitlabId := domain.UserId(userInfo.ID)
	account = &domain.Account{
		Username:     userInfo.Username,
		Name:         userInfo.Name,
		GitlabID:     &gitlabId,
		AuthProvider: domain.AuthProviderGitLab,
		AvatarURL:    userInfo.AvatarURL,
		IsActive:     true,
		JwtSalt:      salt,
	}

	return s.accountRepo.Create(account)
}
