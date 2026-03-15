package models

import "time"

type GitlabTokenId uint

type GitlabToken struct {
	Id           GitlabTokenId
	AccountId    AccountId
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}
