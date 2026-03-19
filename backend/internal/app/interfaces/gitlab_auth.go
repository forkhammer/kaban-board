package interfaces

import (
	domain "main/internal/domain/models"
)

type GitLabAuthServiceInterface interface {
	GetAuthorizationURL() string
	HandleCallback(code string) (*domain.Account, error)
	IsEnabled() bool
	GetValidToken(accountId domain.AccountId) (*domain.GitlabToken, error)
}

type GitLabOAuthToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    int    `json:"created_at"`
}

type GitLabUserInfo struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
}
