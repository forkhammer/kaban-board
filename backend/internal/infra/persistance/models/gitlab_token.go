package models

import "time"

type GitlabToken struct {
	Id           uint      `gorm:"primarykey"`
	AccountId    uint      `gorm:"account_id;not null;uniqueIndex"`
	Account      Account   `gorm:"foreignKey:AccountId"`
	AccessToken  string    `gorm:"access_token;not null"`
	RefreshToken string    `gorm:"refresh_token"`
	ExpiresAt    time.Time `gorm:"expires_at;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
