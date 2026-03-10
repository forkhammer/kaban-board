package models

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	Id           uint `gorm:"id;primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Username     string         `gorm:"username;not null;unique"`
	Password     string         `gorm:"password"`
	Name         string         `gorm:"name"`
	IsActive     bool           `gorm:"is_active;not null"`
	GitlabID     *uint          `gorm:"gitlab_id;unique"`
	AuthProvider string         `gorm:"auth_provider;not null;default:password"`
	AvatarURL    string         `gorm:"avatar_url"`
	Role         string         `gorm:"role;not null;default:viewer"`
}
